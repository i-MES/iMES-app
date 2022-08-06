package utils

import (
	"bytes"
	"context"
	"fmt"
	"hash/crc32"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

func GetAppPath() string {
	// wd, _ := os.Getwd()
	// f, _ := exec.LookPath(os.Args[0])
	_, fileStr, _, _ := runtime.Caller(0)
	approot, _ := filepath.Abs(filepath.Dir(fileStr) + "/../..")
	return approot
}

// .config: 存放最好不要删除的数据
// .cache:  存放随时可以删除的数据
func GetUserCacheDefaultPath() string {
	if home, err := Home(); err == nil {
		dir := home + "/.cache/imes-app/"
		if _, err := os.Stat(dir); err != nil {
			os.MkdirAll(dir, 0750)
		}
		return dir
	} else {
		return ""
	}
}

var _platform string
var _arch string

func InitPlatform(platform string, arch string) {
	_platform = platform
	_arch = arch
}
func GetProcessIdGet() int {
	if _platform == "linux" {
		// return unix.Getpid()
		return os.Getpid()
	}
	if _platform == "windows" {
		return 0
	}
	return 0
}

func GetThreadId() int {
	if _platform == "linux" {
		// return unix.Gettid()
		return 0
	}
	if _platform == "windows" {
		// var user32 *syscall.DLL
		// var GetCurrentThreadId *syscall.Proc
		// var err error

		// user32, err = syscall.LoadDLL("Kernel32.dll")
		// if err != nil {
		// 	log.Debug().Msgf("syscall.LoadDLL fail: %v", err.Error())
		// 	return 0
		// }
		// GetCurrentThreadId, err = user32.FindProc("GetCurrentThreadId")
		// if err != nil {
		// 	log.Debug().Msgf("user32.FindProc fail: %v", err.Error())
		// 	return 0
		// }

		// var pid uintptr
		// pid, _, err = GetCurrentThreadId.Call()

		// return int(pid)
		return 0
	}

	return 0
}

// 遍历查找符合正则的文件
func getAllFile_Walk(root, pattern string) ([]string, error) {
	var result []string
	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
			// 过滤文件夹，such as: __pycache__
			// if _med, _err := filepath.Match("[_|.]*", info.Name()); _err == nil && _med {
			// 	return nil
			// }
			// Walk 函数自己会递归,不需要用户递归，否则子文件夹中会重复找到
			// if recursion {
			// 	if temps, err := getAllFile_Walk(root+"/"+info.Name(), pattern, recursion); err == nil {
			// 		result = append(result, temps...)
			// 	}
			// }
		} else {
			if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
				return err
			} else if matched {
				result = append(result, path)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getAllFile_IOUtil(root, pattern string, recursion bool) ([]string, error) {
	result := []string{}

	// func ReadDir(dirname string) ([]fs.FileInfo, error)
	infos, err := ioutil.ReadDir(root)
	if err != nil {
		log.Debug().Msgf("读取文件目录失败, root=%v, err=%v \n", root, err)
		return result, err
	}

	for _, info := range infos {
		if info.IsDir() {
			if temp, err := getAllFile_IOUtil(root+"/"+info.Name(), pattern, recursion); err != nil {
				return result, err
			} else {
				result = append(result, temp...)
			}
		} else {
			if matched, err := filepath.Match(pattern, filepath.Base(info.Name())); matched && err == nil {
				result = append(result, info.Name())
			} else {
				return nil, err
			}
		}
	}
	return result, nil
}

// 获取指定目录下满足匹配规则的文件
// pattern
// recursion: 是否递归
func GetAllFile(root, pattern string, recursion bool) ([]string, error) {
	if runtime.GOOS == "windows" {
		return getAllFile_IOUtil(root, pattern, recursion)
	} else {
		return getAllFile_Walk(root, pattern)
	}
}

func UUID() string {
	if id, err := uuid.NewUUID(); err == nil {
		return id.String()
	} else {
		return ""
	}
}

// 通过对话框 UI 得到用户选择
func SelectFolder(ctx *context.Context, title string) string {
	if title == "" {
		title = "Open Config Folder"
	}
	_opt := wails.OpenDialogOptions{
		DefaultDirectory: GetAppPath(),
		Title:            title,
	}
	selectedFolder, err := wails.OpenDirectoryDialog(*ctx, _opt)
	if err != nil {
		log.Error().Stack().Err(errors.Wrap(err, "Error on folder opening")).Send()
	}
	return selectedFolder
}

func SelectFile(ctx *context.Context, title, filePattern string) string {
	if title == "" {
		title = "Open Config File"
	}
	if filePattern == "" {
		filePattern = "*.json"
	}
	_opt := wails.OpenDialogOptions{
		DefaultDirectory: "./",
		Title:            title,
		Filters:          []wails.FileFilter{{DisplayName: "File Filter", Pattern: filePattern}},
	}
	selectedFile, err := wails.OpenFileDialog(*ctx, _opt)
	if err != nil {
		log.Error().Stack().Err(errors.Wrap(err, "Error on file opening")).Send()
	}
	return selectedFile
}

// 生成字符串的 hash 值
func Hash(s string) (string, error) {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return fmt.Sprintf("%d", v), nil
	}
	if -v >= 0 {
		return fmt.Sprintf("%d", -v), nil
	}
	return "", fmt.Errorf("generate a hash failed")
}

func Home() (string, error) {
	user, err := user.Current()
	if nil == err {
		return user.HomeDir, nil
	}

	// cross compile support

	if "windows" == runtime.GOOS {
		return homeWindows()
	}

	// Unix-like system, so just assume Unix
	return homeUnix()
}

func homeUnix() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

	// If that fails, try the shell
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "", errors.New("blank output when reading home directory")
	}

	return result, nil
}

func homeWindows() (string, error) {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
	}

	return home, nil
}
