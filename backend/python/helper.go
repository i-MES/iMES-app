package python

/*
#cgo pkg-config: python3-embed
#define PY_SSIZE_T_CLEAN
#include <Python.h>
*/
import "C"

import (
	"fmt"
	"regexp"
)

type Info struct {
	version string
}

func Version_Check(ver string) bool {
	_ver := Py_GetVersion()
	matched, err := regexp.Match(ver, []byte(_ver))
	if err == nil {
		return matched
	}
	return false
}

func PyErr_Occurred() string {
	errstr := ""
	if err := C.PyErr_Occurred(); togo(err) != nil {
		fmt.Println("------ python exception ------")
		C.PyErr_Print()
		fmt.Println("------")
		fmt.Println("sys.path: ", PyImport_GetModule("sys").GetAttrString("path").Repr())
		fmt.Println("------ end ------")
	}
	return errstr
}
