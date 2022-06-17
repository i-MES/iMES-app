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
	"path/filepath"
	"strings"
	"unsafe"
)

func PyImport_ImportModule(name string) *PyObject {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	/*
		C.PyImport_ImportModule:、
		这是个传递函数，接收方需要负责及时 Py_DECREF
		name: 必须是绝对路径：package.submodule.module
		      不能是相对路径：..package.submodule.module
		Return: a new reference to the imported module,所以接收方不用是需 Py_DECREF
	*/
	return togo(C.PyImport_ImportModule(cname))
}

// sys.modules 中有则取出并返回，没有则新建一个空 module 返回
// 空 module 会导致后续无法使用，请确定是否真的想得到空 module 然后自己初始化
func PyImport_AddModule(name string) *PyObject {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	/*
		C.PyImport_AddModule:
		这是个借用函数，接收方不要 Py_DECREF，或者 Py_INCREF 和 Py_DECREF 成对使用。
		首先检查 sys.modules 字典中有没有，有则拿出；
		没有则新建一个 empty 的(仅有 '__name__', '__doc__', '__package__', '__loader__', '__spec__')，并加入 sys.modules 字典；
		出错时返回 NULL；
		该函数不会 load or import module，请用 PyImport_ImportModule()。
	*/
	return togo(C.PyImport_AddModule(cname))
}

// 从 sys.modules 中取出返回，没有则返回 nil
func PyImport_GetModule(name string) *PyObject {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	uname := C.PyUnicode_FromString(cname)
	// defer C.free(unsafe.Pointer(uname))
	// uname 将会是 module 的一部分，
	// free 了 uname，会连带 module 也 free 掉

	/*
		C.PyImport_GetModule
			返回已经 import 的 module，没有则返回 NULL
	*/
	return togo(C.PyImport_GetModule(uname))
}
func PyImport_GetModuleDict() *PyObject {
	return togo(C.PyImport_GetModuleDict())
}

var PySysPathes []string

/*
	导入一个绝对路径描述的 py 文件，并返回 模块名（即：文件名）
	本函数是个传递函数：与 PyImport_ImportModule 保持一致
*/
func PyImport_ImportFile(filename string) *PyObject {
	// step0. 检查是否为绝对路径
	if !filepath.IsAbs(filename) {
		_filename, err := filepath.Abs(filename)
		if err == nil {
			filename = _filename
		}
	}
	// step1. add sys.path
	dir := path.Dir(filename)
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
		PySys_AppendPath(dir)
		PySysPathes = append(PySysPathes, dir)
	} else {
		fmt.Println("Path has added sys.path")
	}
	// step2. load and import module
	fwithsuf := path.Base(filename)
	suf := path.Ext(fwithsuf)
	f := strings.TrimSuffix(fwithsuf, suf)
	_mod := PyImport_GetModule(f)
	// defer _mod.DecRef() 不能 DecRef，否则调用者就 emo 了。
	if _mod == nil {
		fmt.Println("Load and import new module: ", f)
		// way 1:
		// C.PyRun_SimpleString(C.CString(fmt.Sprintf("import %s", f)))
		// way 2:
		// _mod = togo(C.PyImport_Import(C.PyUnicode_FromString(C.CString(f))))
		// way 3:
		_mod = PyImport_ImportModule(f)
		// way 4:
		// _mod = PyImport_AddModule(f) // Error! 只会添加一个 empty 的
	}
	if _mod == nil {
		PyErr_Print()
		fmt.Println(filename)
	}
	// return PyImport_GetModule(f)
	return _mod
}
