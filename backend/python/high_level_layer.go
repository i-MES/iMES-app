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

func PyRun_AnyFile(filename string) (int, error) {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))

	mode := C.CString("r")
	defer C.free(unsafe.Pointer(mode))

	cfile, err := C.fopen(cfilename, mode)
	if err != nil {
		return -1, fmt.Errorf("fail to open '%s': %s", filename, err)
	}
	defer C.fclose(cfile)

	// C.PyRun_AnyFile is a macro, using C.PyRun_AnyFileFlags instead
	return int(C.PyRun_AnyFileFlags(cfile, cfilename, nil)), nil
}

//PyRun_SimpleString : https://docs.python.org/3/c-api/veryhigh.html?highlight=pycompilerflags#c.PyRun_SimpleString
func PyRun_SimpleString(command string) int {
	ccommand := C.CString(command)
	defer C.free(unsafe.Pointer(ccommand))

	// C.PyRun_SimpleString is a macro, using C.PyRun_SimpleStringFlags instead
	return int(C.PyRun_SimpleStringFlags(ccommand, nil))
}
