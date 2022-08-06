package utils

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/pkg/errors"
)

/*
fork from github.com/easyCZ/logrotate
modify:
- 创建时不使用 stdlib 中 log
*/
const (
	lowercase = "abcdefghijklmnopqrstuvwxyz"
	uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits    = "0123456789"
	chars     = lowercase + uppercase + digits
)

const (
	letterIdxBits = 6                        // 6 bits to represent a letter index
	letterIdxMask = (1 << letterIdxBits) - 1 // All 1-bits, as many as letterIdxBits : 0b00111111(63)
)

// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func RandomHash(length int) string {
	b := make([]byte, length)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; {
		if idx := int(rand.Int63() & letterIdxMask); idx < len(chars) {
			b[i] = chars[idx]
			i++
		}
	}
	return string(b)
}

// 生成文件名：2020-03-28_15-00-945-<random-hash>.log
func DefaultFilenameFunc() string {
	return fmt.Sprintf("%s_%s.log", time.Now().UTC().Local().Format(time.RFC3339), RandomHash(3))
}

// Options define configuration options for Writer
type RotateOptions struct {
	// log 存储路径，文件夹不存在会新建
	Directory string

	// 文件最大 size，超过则不再写入，并不会 split 新文件
	// 0：不限制
	MaximumFileSize int64

	// rotate 翻转周期
	// 0：不因时间而翻转
	MaximumLifetime time.Duration

	// 生成文件名函数
	// nil: 使用 DefaultFilenameFunc
	FileNameFunc func() string
}

// Writer is a concurrency-safe writer with file rotation.
type Writer struct {
	// opts are the configuration options for this Writer
	opts RotateOptions

	// f is the currently open file used for appends.
	// Writes to f are only synchronized once Close() is called,
	// or when files are being rotated.
	f *os.File
	// bw is a buffered writer for writing to f
	bw *bufio.Writer
	// bytesWritten is the number of bytes written to f so far,
	// used for size based rotation
	bytesWritten int64
	// ts is the creation timestamp of f,
	// used for time based log rotation
	ts time.Time

	// queue of entries awaiting to be written
	queue chan []byte
	// synchronize write which have started but not been queued up
	pending sync.WaitGroup
	// singal the writer should close
	closing chan struct{}
	// signal the writer has finished writing all queued up entries.
	done chan struct{}
}

// Write writes p into the current file, rotating if necessary.
// Write is non-blocking, if the writer's queue is not full.
// Write is blocking otherwise.
func (w *Writer) Write(p []byte) (n int, err error) {
	select {
	case <-w.closing:
		return 0, errors.Wrap(err, "writer is closing")
	default:
		w.pending.Add(1)
		defer w.pending.Done()
	}

	w.queue <- p
	// 直接写入 channel，并发量太大时，Listen 协程中取会有混乱（截断、穿插）
	// 此处直接写入文件 OK —— 可以证明。
	// if _, err := w.f.Write(p); err != nil {
	// 	fmt.Println("Failed to write to file.", err)
	// }
	// 此处加入 200ms 延时可以解决这个问题
	time.Sleep(time.Microsecond * 200)
	return len(p), nil
}

// Close closes the writer.
// Any accepted writes will be flushed. Any new writes will be rejected.
// Once Close() exits, files are synchronized to disk.
func (w *Writer) Close() error {
	close(w.closing)
	w.pending.Wait()

	close(w.queue)
	<-w.done

	if w.f != nil {
		if err := w.closeCurrentFile(); err != nil {
			return err
		}
	}

	return nil
}

func (w *Writer) listen() {
	for b := range w.queue {
		if w.f == nil {
			if err := w.rotate(); err != nil {
				fmt.Println("Failed to create log file", err)
			}
		}

		size := int64(len(b))

		if w.opts.MaximumFileSize != 0 && size > w.opts.MaximumFileSize {
			fmt.Println("Attempting to write more bytes than allowed by MaximumFileSize. Skipping.")
			continue
		}
		if w.opts.MaximumFileSize != 0 && w.bytesWritten+size > w.opts.MaximumFileSize {
			if err := w.rotate(); err != nil {
				fmt.Println("Failed to rotate log file", err)
			}
		}

		if w.opts.MaximumLifetime != 0 && time.Now().After(w.ts.Add(w.opts.MaximumLifetime)) {
			if err := w.rotate(); err != nil {
				fmt.Println("Failed to rotate log file", err)
			}
		}
		if _, err := w.bw.Write(b); err != nil {
			fmt.Println("Failed to write to file.", err)
		}
		w.bytesWritten += size
		// fmt.Println(size, w.bw.Buffered(), w.bytesWritten)
	}

	close(w.done)
}

func (w *Writer) closeCurrentFile() error {
	if err := w.bw.Flush(); err != nil {
		return errors.Wrap(err, "failed to flush buffered writer")
	}

	if err := w.f.Sync(); err != nil {
		return errors.Wrap(err, "failed to sync current log file")
	}

	if err := w.f.Close(); err != nil {
		return errors.Wrap(err, "failed to close current log file")
	}

	w.bytesWritten = 0
	return nil
}

func (w *Writer) rotate() error {
	if w.f != nil {
		if err := w.closeCurrentFile(); err != nil {
			return err
		}
	}

	path := filepath.Join(w.opts.Directory, w.opts.FileNameFunc())
	f, err := newFile(path)
	if err != nil {
		return errors.Wrapf(err, "failed to create new file at %v", path)
	}

	w.bw = bufio.NewWriter(f)
	w.f = f
	w.bytesWritten = 0
	w.ts = time.Now().UTC()

	return nil
}

// New creates a new concurrency safe Writer which performs log rotation.
func NewRotate(opts RotateOptions) (*Writer, error) {
	if _, err := os.Stat(opts.Directory); os.IsNotExist(err) {
		if err := os.MkdirAll(opts.Directory, os.ModePerm); err != nil {
			return nil, errors.Wrapf(err, "directory %v does not exist and could not be created", opts.Directory)
		}
	}

	if opts.FileNameFunc == nil {
		opts.FileNameFunc = DefaultFilenameFunc
	}

	w := &Writer{
		opts:    opts,
		queue:   make(chan []byte, 1024),
		closing: make(chan struct{}),
		done:    make(chan struct{}),
	}

	w.rotate()
	go w.listen()

	return w, nil
}

func newFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
}
