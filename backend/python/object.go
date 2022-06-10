package python

/*
#cgo pkg-config: python3-embed
#define PY_SSIZE_T_CLEAN
#include <Python.h>
int Cgo_PyDict_Check(PyObject *o){
	return PyDict_Check(o);
}
PyObject* Cgo_PyObject_CallMethod(PyObject *o, const char *name, const char *format){
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

// PyObject Creater
func PyUnicode_FromString(u string) *PyObject {
	cu := C.CString(u)
	defer C.free(unsafe.Pointer(cu))

	return togo(C.PyUnicode_FromString(cu))
}

// PyObject 方法
func (pyObject *PyObject) IncRef() {
	C.Py_IncRef(toc(pyObject))
}

func (pyObject *PyObject) DecRef() {
	C.Py_DecRef(toc(pyObject))
}

func (pyObject *PyObject) HasAttr(attr_name *PyObject) bool {
	return C.PyObject_HasAttr(toc(pyObject), toc(attr_name)) == 1
}

func (pyObject *PyObject) HasAttrString(attr_name string) bool {
	cattr_name := C.CString(attr_name)
	defer C.free(unsafe.Pointer(cattr_name))

	return C.PyObject_HasAttrString(toc(pyObject), cattr_name) == 1
}

func (pyObject *PyObject) GetAttr(attr_name *PyObject) *PyObject {
	return togo(C.PyObject_GetAttr(toc(pyObject), toc(attr_name)))
}

func (pyObject *PyObject) GetAttrString(attr_name string) *PyObject {
	cattr_name := C.CString(attr_name)
	defer C.free(unsafe.Pointer(cattr_name))

	return togo(C.PyObject_GetAttrString(toc(pyObject), cattr_name))
}

func (pyObject *PyObject) CallObject(args *PyObject) *PyObject {
	return togo(C.PyObject_CallObject(toc(pyObject), toc(args)))
}

// Call a  method without args
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

// 类型识别与数据提取
// Unicode 类型的 PyObject 才能调用的 UTF8 转换函数
func (pyObject *PyObject) UTF8() string {
	_cutf8 := C.PyUnicode_AsUTF8(toc(pyObject))
	return C.GoString(_cutf8)
}

// 获取 PyObject 的字符串表达: 类型（值）
func (pyObject *PyObject) Repr() string {
	_str := C.PyObject_Repr(toc(pyObject))
	defer C.Py_DecRef(_str)
	if _str == nil {
		return ""
	} else {
		return togo(_str).UTF8()
	}
}

// 获取 PyObject 的字符串表达，与 str(object) 等效，只有值
func (pyObject *PyObject) Str() string {
	_str := C.PyObject_Str(toc(pyObject))
	defer C.Py_DecRef(_str)
	if _str == nil {
		return ""
	} else {
		return togo(_str).UTF8()
	}
}

// 获取 PyObject 的类型并转换成字符串，与 print(type(object)) 等效
func (pyObject *PyObject) Type() string {
	_type := C.PyObject_Type(toc(pyObject))
	defer C.Py_DecRef(_type)
	return togo(_type).Repr()
}

// func (pyObject *PyObject) Bytes() int {
// 	_bytes := C.PyObject_Bytes(toc(pyObject))
// 	defer C.Py_DecRef(_bytes)
// 	fmt.Println(_bytes)
// 	fmt.Println(C.PyBytes_Size(_bytes))
// 	if _bytes == nil {
// 		return 0
// 	} else {
// 		return int(C.PyBytes_Size(_bytes))
// 	}
// }

// 只有数值型 PyObject 才能使用的获取数值函数
func (pyObject *PyObject) Number() int {
	// a := *(*int)(unsafe.Pointer(o.ob_type))
	return int(C.PyLong_AsLong(toc(pyObject)))
}

// 只有数值型 PyObject 才能使用的获取 Long 函数（类似 Long 是 Number 的子类）
func (pyObject *PyObject) Long() int {
	// a := *(*int)(unsafe.Pointer(o.ob_type))
	return int(C.PyLong_AsLong(toc(pyObject)))
}
