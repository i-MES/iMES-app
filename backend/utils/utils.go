package utils

// #include <unistd.h>
import "C"

import (
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime"

	"github.com/google/uuid"
	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

func GetAppPath() string {
	// wd, _ := os.Getwd()
	// f, _ := exec.LookPath(os.Args[0])
	_, fileStr, _, _ := runtime.Caller(0)
	approot, _ := filepath.Abs(filepath.Dir(fileStr) + "/../..")
	// fmt.Println("AppRoot: ", approot)
	return approot
}

var _platform string
var _arch string

func InitPlatform(platform string, arch string) {
	_platform = platform
	_arch = arch
}
func GetProcessId() int {
	if _platform == "linux" {
		return int(C.getpid())
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
		// 	fmt.Printf("syscall.LoadDLL fail: %v", err.Error())
		// 	return 0
		// }
		// GetCurrentThreadId, err = user32.FindProc("GetCurrentThreadId")
		// if err != nil {
		// 	fmt.Printf("user32.FindProc fail: %v", err.Error())
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
func getAllFile_Walk(root, pattern string, recursion bool) ([]string, error) {
	var result []string
	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if recursion {
				if temps, err := getAllFile_Walk(root+"/"+info.Name(), pattern, recursion); err == nil {
					result = append(result, temps...)
				}
			}
		} else if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			result = append(result, path)
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
		fmt.Printf("读取文件目录失败，root=%v, err=%v \n", root, err)
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
		return getAllFile_Walk(root, pattern, recursion)
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
func OpenFolder(ctx *context.Context, title string) string {
	if title == "" {
		title = "Open Config Folder"
	}
	_opt := wails.OpenDialogOptions{
		DefaultDirectory: GetAppPath(),
		Title:            title,
	}
	selectedFolder, err := wails.OpenDirectoryDialog(*ctx, _opt)
	if err != nil {
		log.Panic("Error on folder opening", err.Error())
	}
	return selectedFolder
}

func OpenFile(ctx *context.Context, title, filePattern string) string {
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
		log.Panic("Error on file opening", err.Error())
	}
	return selectedFile
}
