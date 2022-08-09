package python

/*
#cgo pkg-config: python3-embed
#define PY_SSIZE_T_CLEAN
#include <Python.h>
PyObject* Cgo_PyObject_CallMethod(PyObject *o, const char *name, const char *format){
	// PyObject_CallMethod:
	// 		调用 obj 对象中名为 name 的方法并附带可变数量的 C 参数。
	// 		这些 C 参数由 Py_BuildValue() 格式字符串来描述并应当生成一个元组。
	// 		如果参数都是 PyObject*，没有 C 型的，则 PyObject_CallMethodObjArgs() 是更快速的选择。
	return PyObject_CallMethod(o, name, format);
}
PyObject* Cgo_PyObject_CallMethodObjArgs(PyObject *obj, PyObject *name, int argc, PyObject **argv) {
    PyObject *result = NULL;
    switch (argc) {
        case 0:
            return PyObject_CallMethodObjArgs(obj, name, NULL);
        case 1:
            return PyObject_CallMethodObjArgs(obj, name, argv[0], NULL);
        case 2:
            return PyObject_CallMethodObjArgs(obj, name, argv[0], argv[1], NULL);
        case 3:
            return PyObject_CallMethodObjArgs(obj, name, argv[0], argv[1], argv[2], NULL);
    }
    return result;
}
*/
import "C"
import (
	"unsafe"
)

/*
Class 类型的 PyObject 相关函数
*/

func (pyObject *PyObject) GetAttr(attr_name *PyObject) *PyObject {
	return togo(C.PyObject_GetAttr(toc(pyObject), toc(attr_name)))
}

func (pyObject *PyObject) GetAttrString(attr_name string) *PyObject {
	cattr_name := C.CString(attr_name)
	defer C.free(unsafe.Pointer(cattr_name))

	return togo(C.PyObject_GetAttrString(toc(pyObject), cattr_name))
}

/*
 Call a  method without args
 可用于实例化 Class、调用 Function()、调用 object.method()
*/
func (pyObject *PyObject) CallMethod(methodName string) *PyObject {
	return togo(C.Cgo_PyObject_CallMethod(toc(pyObject), C.CString(methodName), nil))
}

// Call a  method with args(max len = 3)
func (pyObject *PyObject) CallMethodArgs(methodName string, args ...*PyObject) *PyObject {
	mname := PyUnicode_FromString(methodName)
	defer mname.DecRef()

	if len(args) > 3 {
		panic("CallMethodObjArgs: too many arguments")
	}
	if len(args) == 0 {
		return togo(C.Cgo_PyObject_CallMethodObjArgs(toc(pyObject), toc(mname), 0, (**C.PyObject)(nil)))
	}

	cargs := make([]*C.PyObject, len(args))
	for i, arg := range args {
		cargs[i] = toc(arg)
	}
	return togo(C.Cgo_PyObject_CallMethodObjArgs(toc(pyObject), toc(mname), C.int(len(args)), (**C.PyObject)(unsafe.Pointer(&cargs[0]))))
}
