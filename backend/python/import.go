package python

/*
#cgo pkg-config: python3-embed
#define PY_SSIZE_T_CLEAN
#include <Python.h>
*/
import "C"
import (
	"fmt"
	"path"
	"strings"
	"unsafe"
)

func PyImport_ImportModule(name string) *PyObject {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	return togo(C.PyImport_ImportModule(cname))
}
func PyImport_ImportModuleEx(name string, globals, locals, fromlist *PyObject) *PyObject {
	return PyImport_ImportModuleLevel(name, globals, locals, fromlist, 0)
}

// Level: 0 (默认值) -仅执行绝对导入(x.y); 其他-相对导入(..x.y)
func PyImport_ImportModuleLevel(name string, globals, locals, fromlist *PyObject, level int) *PyObject {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	return togo(C.PyImport_ImportModuleLevel(cname, toc(globals), toc(locals), toc(fromlist), C.int(level)))
}
func PyImport_GetModuleDict() *PyObject {
	return togo(C.PyImport_GetModuleDict())
}
func PyImport_GetModule(name string) *PyObject {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	uname := C.PyUnicode_FromString(cname)
	// defer C.free(unsafe.Pointer(uname))

	return togo(C.PyImport_GetModule(uname))
}
func PyImport_AddModule(name string) *PyObject {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	return togo(C.PyImport_AddModule(cname))
}

var PySysPathes []string

// 导入一个绝对路径描述的 py 文件，并返回 模块名（即：文件名）
func PyImport_ImportFile(filepath string) *PyObject {
	// step1. add sys.path
	dir := path.Dir(filepath)
	if dir != "." {
		if PySysPathes == nil {
			PySysPathes = make([]string, 0)
		}
		_pathadded := false
		for _, d := range PySysPathes {
			if d == dir {
				_pathadded = true
				break
			}
		}
		if !_pathadded {
			fmt.Println("Add a new path to sys.path")
			cmd := fmt.Sprintf("sys.path.append('%s')", dir)
			C.PyRun_SimpleString(C.CString(cmd))
			PySysPathes = append(PySysPathes, dir)
		} else {
			fmt.Println("Path has added sys.path")
		}
	}
	// step2. load and import module
	fwithsuf := path.Base(filepath)
	suf := path.Ext(fwithsuf)
	f := strings.TrimSuffix(fwithsuf, suf)
	_mod := PyImport_GetModule(f)
	// defer _mod.DecRef() 不能 DecRef，否则调用者就 emo 了。
	if _mod == nil {
		fmt.Println("Load and import new module: ", f)
		// C.PyImport_Import(C.PyUnicode_FromString(C.CString(f))) // 与下面一行同效
		// C.PyImport_ImportModule(C.CString(f))
		// C.PyRun_SimpleString(C.CString(fmt.Sprintf("import %s", f)))
		_mod = PyImport_AddModule(f)
	} else {
		fmt.Println("Module has loaded")
	}

	// return PyImport_AddModule(f)
	return _mod
}
