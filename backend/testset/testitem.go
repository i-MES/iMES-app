package testset

/*
#cgo pkg-config: python-3.9-embed
#define PY_SSIZE_T_CLEAN
#include <Python.h>
*/
import "C"
import (
	"fmt"
)

// 测试项
type TestItem struct {
	Title    string `json:"title"`
	Desc     string `json:"desc"`
	Filename string `json:"filename"`
	Funcname string `json:"funcname"`
	Sequence int    `json:"sequence"`
}

func (ti *TestItem) Run() {
	// runtime.LockOSThread()
	// _gstate := C.PyGILState_Ensure()
	// C.PyGILState_Release(_gstate)
	fmt.Println(C.GoString(C.Py_GetVersion()))
	fmt.Println(ti)

	pModule := C.PyImport_ImportModule(C.CString(ti.Filename))
	if pModule != nil {
		pFunc := C.PyObject_GetAttrString(pModule, C.CString(ti.Funcname))
		if pFunc != nil {
			pValue := C.PyObject_CallObject(pFunc, nil)
			if pValue != nil {
				fmt.Println(pValue)
			}
		}
	}

	fmt.Println("Run TI: ", ti.Title, ti.Desc, ti.Funcname, ti.Sequence)
}
