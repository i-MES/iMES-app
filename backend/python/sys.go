package python

/*
#cgo pkg-config: python3-embed
#define PY_SSIZE_T_CLEAN
#include <Python.h>
*/
import "C"

func PySys_GetSysPath() string {
	return PyImport_GetModule("sys").GetAttrString("path").Str()
}
func PySys_GetSysModules() string {
	return PyImport_GetModule("sys").GetAttrString("modules").Str()
}

// Append path to sys.path
func PySys_AppendSysPath(path string) error {
	// way 1:
	// cmd := fmt.Sprintf("sys.path.append('%s')", path)
	// C.PyRun_SimpleString(C.CString(cmd))

	// way 2:
	args := C.PyUnicode_FromString(C.CString(path))
	defer C.Py_DecRef(args)
	PyImport_GetModule("sys").GetAttrString("path").CallMethodArgs("append", togo(args))

	// way 3:
	// cpath := C.CString(path)
	// defer C.free(unsafe.Pointer(cpath))
	// wpath := C.Py_DecodeLocale(cpath, nil)
	// if wpath == nil {
	// 	return fmt.Errorf("fail to call Py_DecodeLocale on '%s'", path)
	// }
	// defer C.PyMem_RawFree(unsafe.Pointer(wpath))
	// C.PySys_SetPath(wpath)

	return nil
}

// Append path to sys.path
func PySys_RemoveSysPath(path string) error {
	args := C.PyUnicode_FromString(C.CString(path))
	defer C.Py_DecRef(args)
	PyImport_GetModule("sys").GetAttrString("path").CallMethodArgs("remove", togo(args))
	return nil
}
