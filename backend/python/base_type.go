package python

/*
#cgo pkg-config: python3-embed
#define PY_SSIZE_T_CLEAN
#include <Python.h>
*/
import "C"
import "unsafe"

func PyDict_GetItemString(p *PyObject, key string) *PyObject {
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))

	return togo(C.PyDict_GetItemString(toc(p), ckey))
}

func PyList_New(len int) *PyObject {
	return togo(C.PyList_New(C.Py_ssize_t(len)))
}

func PyTuple_New(len int) *PyObject {
	return togo(C.PyTuple_New(C.Py_ssize_t(len)))
}

func PyTuple_SetItem(p *PyObject, pos int, o *PyObject) int {
	return int(C.PyTuple_SetItem(toc(p), C.Py_ssize_t(pos), toc(o)))
}
