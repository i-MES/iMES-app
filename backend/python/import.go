package python

/*
#cgo pkg-config: python3-embed
#define PY_SSIZE_T_CLEAN
#include <Python.h>
*/
import "C"
import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"unsafe"

	"github.com/rs/zerolog/log"
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
	root: 作为 module path
	file - root: 作为 module name( / 替换为 .)
*/
func PyImport_ImportFile(root, filename string) *PyObject {
	// step0. 检查是否存在、绝对路径
	if _, _err := os.Stat(filename); _err != nil {
		// 文件不存在
		return nil
	}
	if !filepath.IsAbs(root) {
		_root, err := filepath.Abs(root)
		if err == nil {
			root = _root
		}
	}
	if !filepath.IsAbs(filename) {
		_filename, err := filepath.Abs(filename)
		if err == nil {
			filename = _filename
		}
	}

	// step1. add sys.path
	if PySysPathes == nil {
		PySysPathes = make([]string, 0)
	}
	_pathadded := false
	for _, d := range PySysPathes {
		if d == root {
			_pathadded = true
			break
		}
	}
	if !_pathadded {
		log.Debug().Msgf("Add a new path(%s) to sys.path\n", root)
		PySys_AppendSysPath(root)
		PySysPathes = append(PySysPathes, root)
	} else {
		log.Debug().Msgf("Path(%s) has added sys.path\n", root)
	}

	// step2. load and import module
	fwithsuf := strings.ReplaceAll(filename[len([]rune(root))+1:], "/", ".")
	suf := path.Ext(fwithsuf)
	modname := strings.TrimSuffix(fwithsuf, suf) // 去掉 .py 后缀
	_mod := PyImport_GetModule(modname)
	// defer _mod.DecRef() 不能 DecRef，否则调用者就 emo 了。
	if _mod == nil {
		log.Debug().Msgf("Load and import new module: %v", modname)
		// way 1:
		// C.PyRun_SimpleString(C.CString(fmt.Sprintf("import %s", f)))
		// way 2:
		// _mod = togo(C.PyImport_Import(C.PyUnicode_FromString(C.CString(f))))
		// way 3:
		_mod = PyImport_ImportModule(modname)
		// way 4:
		// _mod = PyImport_AddModule(f) // Error! 只会添加一个 empty 的
	}
	if _mod == nil {
		PyErr_Print()
		log.Debug().Msg(filename)
	}

	// step3. 查看新导入的 mod 的路径是否与入参匹配
	// 如：入参是希望导入 x/y/z.py，但已经导入过 m/n/z.py，
	// 则：m/n 一直在 sys.path 中，即使  x/y/z.py 不存在，也能 import z 成功
	modfile := _mod.GetAttrString("__file__").Str()
	if filename != modfile {
		log.Error().Stack().Msg("Import wrong module: " + root + "!=" + modfile)
		PySys_RemoveSysPath(path.Dir(modfile))
		log.Debug().Msg(PySys_GetSysPath())
		return nil
	}

	// return PyImport_GetModule(f)
	return _mod
}

func PyImport_AddPathAndImportModule(modulepath, modulename string) *PyObject {
	// step0. 检查是否为绝对路径
	if !filepath.IsAbs(modulepath) {
		_modulepath, err := filepath.Abs(modulepath)
		if err == nil {
			modulepath = _modulepath
		}
	}
	// step1. add sys.path
	if PySysPathes == nil {
		PySysPathes = make([]string, 0)
	}
	_pathadded := false
	for _, d := range PySysPathes {
		if d == modulepath {
			_pathadded = true
			break
		}
	}
	if !_pathadded {
		log.Debug().Msgf("Add a new path(%s) to sys.path\n", modulepath)
		PySys_AppendSysPath(modulepath)
		PySysPathes = append(PySysPathes, modulepath)
	} else {
		log.Debug().Msgf("Path(%s) has added sys.path\n", modulepath)
	}

	// step2. load and import module
	_mod := PyImport_GetModule(modulename)
	if _mod == nil {
		log.Debug().Msg("Load and import new module: " + modulename)
		_mod = PyImport_ImportModule(modulename)
	}
	if _mod == nil {
		PyErr_Print()
	}

	// return PyImport_GetModule(f)
	return _mod
}
