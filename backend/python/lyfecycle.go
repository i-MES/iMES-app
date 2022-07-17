package python

/*
#cgo pkg-config: python3-embed
#define PY_SSIZE_T_CLEAN
#include <Python.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

var (
	programName *C.wchar_t
	pythonPath  *C.wchar_t
	pythonHome  *C.wchar_t
)

func Py_IsInitialized() bool {
	return C.Py_IsInitialized() != 0
}

func Py_Initialize() {
	C.Py_Initialize()
}

// syspath: append to sys.path
func Py_InitializeFromConfig(syspath []string) {
	var config C.PyConfig
	// var wc *C.wchar_t
	C.PyConfig_InitPythonConfig(&config)
	for _, _sp := range syspath {
		sp := C.CString(_sp)
		wsp := C.Py_DecodeLocale(sp, nil)
		C.PyWideStringList_Append(&config.module_search_paths, wsp)
		C.free(unsafe.Pointer(sp))
		C.free(unsafe.Pointer(wsp))
	}
	C.Py_InitializeFromConfig(&config)
}

func Py_Finalize() {
	C.Py_Finalize()
}

func Py_GetVersion() string {
	cversion := C.Py_GetVersion()
	return C.GoString(cversion)
}

// 设置程序名，类似命令行执行时 main(argc, *argv) 中的 argv[0]
//
// 只能在 Py_Initialize 之前 Set，运行中不能修改。
// 会参与 Prefix 的生成，如：
//     if the program name is '/usr/local/bin/python', the prefix is '/usr/local'.
func Py_SetProgramName(name string) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	// wchar_t *Py_DecodeLocale(const char *arg, size_t *size)
	// 解析字符串为 wchar_t
	// 但该函数不能在 Python preinitialized 之前调用
	newProgramName := C.Py_DecodeLocale(cname, nil)
	if newProgramName == nil {
		return fmt.Errorf("fail to call Py_DecodeLocale on '%s'", name)
	}
	C.Py_SetProgramName(newProgramName)

	//no operation is performed if nil
	C.PyMem_RawFree(unsafe.Pointer(programName))
	programName = newProgramName

	return nil
}
func Py_GetProgramName() (string, error) {
	wcname := C.Py_GetProgramName()
	if wcname == nil {
		return "", nil
	}
	cname := C.Py_EncodeLocale(wcname, nil)
	if cname == nil {
		return "", fmt.Errorf("fail to call Py_EncodeLocale")
	}
	defer C.PyMem_Free(unsafe.Pointer(cname))

	return C.GoString(cname), nil
}

// Set the default module search path.
//
// 会导致清空 sys.prefix 和 sys.exec_prefix
func Py_SetPath(path string) error {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	newPath := C.Py_DecodeLocale(cpath, nil)
	if newPath == nil {
		return fmt.Errorf("fail to call Py_DecodeLocale on '%s'", path)
	}
	C.Py_SetPath(newPath)

	C.PyMem_RawFree(unsafe.Pointer(pythonPath))
	pythonHome = newPath

	return nil
}
func Py_GetPath() (string, error) {
	wcname := C.Py_GetPath()
	if wcname == nil {
		return "", nil
	}
	cname := C.Py_EncodeLocale(wcname, nil)
	if cname == nil {
		return "", fmt.Errorf("fail to call Py_EncodeLocale")
	}
	defer C.PyMem_Free(unsafe.Pointer(cname))

	return C.GoString(cname), nil
}

// Set the location of the standard Python libraries。
//
// 命令行时从 PYTHONHOME 环境变量获取。
// 默认在 prefix/lib/pythonX.Y 和 exec_prefix/lib/pythonX.Y 中搜索。
//
// venv 虚拟环境不包含标准库，只包含 site-packages 3rd 库，所以不要配置成 venv 路径了。
func Py_SetPythonHome(home string) error {
	chome := C.CString(home)
	defer C.free(unsafe.Pointer(chome))

	newHome := C.Py_DecodeLocale(chome, nil)
	if newHome == nil {
		return fmt.Errorf("fail to call Py_DecodeLocale on '%s'", home)
	}
	C.Py_SetPythonHome(newHome)

	C.PyMem_RawFree(unsafe.Pointer(pythonHome))
	pythonHome = newHome

	return nil
}
func Py_GetPythonHome() (string, error) {
	wchome := C.Py_GetPythonHome()
	if wchome == nil {
		return "", nil
	}
	chome := C.Py_EncodeLocale(wchome, nil)
	if chome == nil {
		return "", fmt.Errorf("fail to call Py_EncodeLocale")
	}
	defer C.PyMem_Free(unsafe.Pointer(chome))

	return C.GoString(chome), nil
}
