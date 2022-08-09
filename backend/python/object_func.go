package python

/*
#cgo pkg-config: python3-embed
#define PY_SSIZE_T_CLEAN
#include <Python.h>
PyObject* Cgo_PyObject_CallFunctionObjArgs(PyObject *callable, int argc, PyObject **argv) {
    PyObject *result = NULL;
    switch (argc) {
        case 0:
            return PyObject_CallMethodObjArgs(callable, NULL);
        case 1:
            return PyObject_CallMethodObjArgs(callable, argv[0], NULL);
        case 2:
            return PyObject_CallMethodObjArgs(callable, argv[0], argv[1], NULL);
        case 3:
            return PyObject_CallMethodObjArgs(callable, argv[0], argv[1], argv[2], NULL);
    }
    return result;
}
*/
import "C"
import (
	"unsafe"
)

/*
functions 类型 PyObject 才能调用的方法
*/

// 调用一个 Python 函数的高层封装
func CallBuiltinsFunc(funcname string) *PyObject {
	builtins := PyEval_GetBuiltins()
	fu := PyDict_GetItemString(builtins, funcname)

	emptyList := PyList_New(0)
	args := PyTuple_New(1)
	defer args.DecRef()
	PyTuple_SetItem(args, 0, emptyList)

	return fu.Call(args, nil) // or fu.CallObject(args)
}

// 待测试 API
func CompileAndRun(filename, funcname string) *PyObject {
	ver := PyUnicode_FromString("__version__ = '3.0'")
	defer ver.DecRef()

	fn := PyUnicode_FromString(filename)
	defer fn.DecRef()

	mode := PyUnicode_FromString("exec")
	defer mode.DecRef()

	builtins := PyEval_GetBuiltins()
	compile := PyDict_GetItemString(builtins, "compile")
	code := compile.CallFunctionObjArgs(ver, fn, mode)
	defer code.DecRef()

	module := PyImport_ExecCodeModule(funcname, code)
	return module
}

/*
从一个 Callable Object 调用，入参即函数入参，入参可以 nil
*/
func (pyObject *PyObject) Call(args *PyObject, kwargs *PyObject) *PyObject {
	return togo(C.PyObject_Call(toc(pyObject), toc(args), toc(kwargs)))
}

func (pyObject *PyObject) CallObject(args *PyObject) *PyObject {
	return togo(C.PyObject_CallObject(toc(pyObject), toc(args)))
}

func (pyObject *PyObject) CallFunctionObjArgs(args ...*PyObject) *PyObject {

	if len(args) > 3 {
		panic("CallFunctionObjArgs: too many arrguments")
	}
	if len(args) == 0 {
		return togo(C.Cgo_PyObject_CallFunctionObjArgs(toc(pyObject), 0, (**C.PyObject)(nil)))
	}

	cargs := make([]*C.PyObject, len(args), len(args))
	for i, arg := range args {
		cargs[i] = toc(arg)
	}
	return togo(C.Cgo_PyObject_CallFunctionObjArgs(toc(pyObject), C.int(len(args)), (**C.PyObject)(unsafe.Pointer(&cargs[0]))))
}
