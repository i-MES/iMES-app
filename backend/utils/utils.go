package utils

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
