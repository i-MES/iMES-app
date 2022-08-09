package python

/*
#cgo pkg-config: python3-embed
#define PY_SSIZE_T_CLEAN
#include <Python.h>
int Cgo_PyDict_Check(PyObject *o){
	return PyDict_Check(o);
}
*/
import "C"
import "unsafe"

/*
所有类型 PyObject 都能使用的函数
*/

type PyObject C.PyObject

// Helper 函数
func togo(cobject *C.PyObject) *PyObject {
	return (*PyObject)(cobject)
}
func toc(object *PyObject) *C.PyObject {
	return (*C.PyObject)(object)
}

// 类型检测
func PyCallable_Check(o *PyObject) bool {
	return C.PyCallable_Check(toc(o)) == 1
}
func PyNumber_Check(o *PyObject) bool {
	return C.PyNumber_Check(toc(o)) != 0
}
func PyLong_Check(o *PyObject) bool {
	return C.PyLong_AsDouble(toc(o)) != 0
}
func PyDict_Check(o *PyObject) bool {
	return C.Cgo_PyDict_Check(toc(o)) != 0
}

// 引用相关
func (pyObject *PyObject) IncRef() {
	C.Py_IncRef(toc(pyObject))
}

func (pyObject *PyObject) DecRef() {
	C.Py_DecRef(toc(pyObject))
}

// 属性查找相关
func (pyObject *PyObject) HasAttr(attr_name *PyObject) bool {
	return C.PyObject_HasAttr(toc(pyObject), toc(attr_name)) == 1
}

func (pyObject *PyObject) HasAttrString(attr_name string) bool {
	cattr_name := C.CString(attr_name)
	defer C.free(unsafe.Pointer(cattr_name))

	return C.PyObject_HasAttrString(toc(pyObject), cattr_name) == 1
}
