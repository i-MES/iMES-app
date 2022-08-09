package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLog(t *testing.T) {
}

func TestLogRotate(t *testing.T) {
	setup := func(t *testing.T) (string, func()) {
		dir, err := ioutil.TempDir("", "") // 在 /tmp/ 下创建临时文件夹
		require.NoError(t, err)

		cleanup := func() {
			t.Logf("Delete temp path: %v", dir)
			require.NoError(t, os.RemoveAll(dir)) // 删除临时文件夹
		}

		return dir, cleanup
	}

	t.Run("creates target directory if it does not exist", func(t *testing.T) {
		dir, cleanup := setup(t)
		defer cleanup()

		dir = filepath.Join(dir, "foo")
		w, err := NewRotate(RotateOptions{
			Directory: dir,
		})
		require.NoError(t, err)
		require.NoError(t, w.Close(), "must close writer")

		// 检查 /tmp/xxx/foo 是否已经创建
		f, err := os.Stat(dir)
		require.NoError(t, err)
		require.True(t, f.IsDir(), "must create directory")
	})

	t.Run("create, write, close", func(t *testing.T) {
		dir, cleanup := setup(t)
		defer cleanup()

		w, err := NewRotate(RotateOptions{
			Directory: dir,
		})
		require.NoError(t, err, "must construct writer")

		// 写入 log 中一些内容
		message := []byte("message")
		_, err = w.Write(message)
		require.NoError(t, err, "write must succeed")
		require.NoError(t, w.Close(), "must close writer")

		files, err := ioutil.ReadDir(dir)
		require.NoError(t, err)

		// 检查是否仅有 1 个 log 文件，并且内容相同
		require.Len(t, files, 1, "must write exactly one file")
		written, err := ioutil.ReadFile(filepath.Join(dir, files[0].Name()))
		require.NoError(t, err, "must read file")
		require.Equal(t, message, written)
	})

	t.Run("rotates on file size", func(t *testing.T) {
		dir, cleanup := setup(t)
		defer cleanup()

		max := 128
		w, err := NewRotate(RotateOptions{
			Directory:       dir,
			MaximumFileSize: int64(max), // 达到 128 个 Byte 开始翻转
		})
		require.NoError(t, err)

		// fill up the first file
		_, err = w.Write([]byte(strings.Repeat("a", max))) // 写入 128 个字符 a
		require.NoError(t, err, "must write")

		// write more, should create a new file
		_, err = w.Write([]byte("b")) // 继续写入一个 b
		require.NoError(t, err)

		require.NoError(t, w.Close())

		files, err := ioutil.ReadDir(dir)
		require.NoError(t, err)

		require.Len(t, files, 2, "must produce 2 files") // 检查是否有 2 个 log 文件
	})

	t.Run("rotates on lifetime", func(t *testing.T) {
		dir, cleanup := setup(t)
		defer cleanup()

		lifetime := time.Second
		w, err := NewRotate(RotateOptions{
			Directory:       dir,
			MaximumLifetime: lifetime, // 1s 就翻转
		})
		require.NoError(t, err)

		// keep writing until lifetime + half of lifetime (middle of ticks) elapses
		end := time.Now().Add(lifetime + lifetime/2)
		for time.Now().Before(end) {
			_, err = w.Write([]byte("message"))
			require.NoError(t, err)
		}

		require.NoError(t, w.Close())
		files, err := ioutil.ReadDir(dir)
		require.NoError(t, err)
		require.Len(t, files, 2, "should produce 2 files")
	})

	// 验证并发写是否内容完整
	t.Run("concurrent writes right number", func(t *testing.T) {
		dir, cleanup := setup(t)
		defer cleanup()

		w, err := NewRotate(RotateOptions{
			Directory: dir,
		})
		require.NoError(t, err)
		rows := 1000
		writers := 10
		messageSize := 10

		var wg sync.WaitGroup
		for i := 0; i < writers; i++ {
			wg.Add(1)
			go func(i int) {
				for j := 0; j < rows; j++ {
					_, err := w.Write([]byte(strings.Repeat(fmt.Sprintf("%d", i), messageSize)))
					require.NoError(t, err)
				}
				wg.Done()
			}(i)
		}

		wg.Wait()
		require.NoError(t, w.Close(), "must close")

		files, err := ioutil.ReadDir(dir)
		require.NoError(t, err)
		require.Len(t, files, 1, "must write a single file")
		t.Log(files[0].Size())
		require.Equal(t, int64(rows*writers*messageSize), files[0].Size(), "must write all bytes")
	})

	// 验证并发写是否内容顺序正确
	t.Run("concurrent writes right queue", func(t *testing.T) {
		dir, cleanup := setup(t)
		defer cleanup()

		w, err := NewRotate(RotateOptions{
			Directory: dir,
		})
		require.NoError(t, err)
		rows := 1000
		writers := 10
		messageSize := 100

		var wg sync.WaitGroup
		for i := 0; i < writers; i++ {
			wg.Add(1)
			go func(i int) {
				for j := 0; j < rows; j++ {
					// 每次写入是一行（末尾 \n）
					_, err := w.Write([]byte(strings.Repeat(fmt.Sprintf("%d", i), messageSize-1) + "\n"))
					require.NoError(t, err)
				}
				wg.Done()
			}(i)
		}

		wg.Wait()
		require.NoError(t, w.Close(), "must close")

		files, err := ioutil.ReadDir(dir)
		require.NoError(t, err)
		require.Len(t, files, 1, "must write a single file")
		require.Equal(t, int64(rows*writers*messageSize), files[0].Size(), "must write all bytes")

		fp := filepath.Join(dir, files[0].Name())
		t.Log(fp)
		bs, e := ioutil.ReadFile(fp)
		require.NoError(t, e)

		// 验证每行内部是否数字相同
		for i := 0; i < writers*rows; i++ {
			for j := 0; j < messageSize-3; j++ {
				k := i*messageSize + j
				require.Equal(t, bs[k], bs[k+1], "must same letter in one message,but find %c(%d)%c(%d) in line %d, char index %d", bs[k], k, bs[k+1], k+1, i, j)
			}
		}
	})
}
