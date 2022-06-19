package utils

// #include <unistd.h>
import "C"

import (
	"path/filepath"
	"runtime"
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
