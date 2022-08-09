package python

/*
#cgo pkg-config: python3-embed
#define PY_SSIZE_T_CLEAN
#include <Python.h>
*/
import "C"

func PyEval_GetBuiltins() *PyObject {
	return togo(C.PyEval_GetBuiltins())
}

func PyEval_GetLocals() *PyObject {
	return togo(C.PyEval_GetLocals())
}

func PyEval_GetGlobals() *PyObject {
	return togo(C.PyEval_GetGlobals())
}

func PyEval_GetFuncName(pyFunc *PyObject) string {
	return C.GoString(C.PyEval_GetFuncName(toc(pyFunc)))
}

func PyEval_GetFuncDesc(pyFunc *PyObject) string {
	return C.GoString(C.PyEval_GetFuncDesc(toc(pyFunc)))
}
