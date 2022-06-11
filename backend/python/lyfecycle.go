package python

/*
#cgo pkg-config: python3-embed
#define PY_SSIZE_T_CLEAN
#include <Python.h>
*/
import "C"

func Py_IsInitialized() bool {
	return C.Py_IsInitialized() != 0
}

func Py_Initialize() {
	C.Py_Initialize()
}

func Py_Finalize() {
	C.Py_Finalize()
}

func Py_GetVersion() string {
	cversion := C.Py_GetVersion()
	return C.GoString(cversion)
}
