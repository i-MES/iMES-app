package utils

// #include <unistd.h>
import "C"

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/google/uuid"
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
func WalkMatch(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil

}

func UUID() string {
	if id, err := uuid.NewUUID(); err == nil {
		return id.String()
	} else {
		return ""
	}
}
