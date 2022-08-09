package python

/*
#cgo pkg-config: python3-embed
#define PY_SSIZE_T_CLEAN
#include <Python.h>
*/
import "C"
import "unsafe"

/*
基础类型（int,string,unicode,...） PyObject 都能使用的函数
*/

// 创建 Unicode 类型 PyObject
func PyUnicode_FromString(u string) *PyObject {
	cu := C.CString(u)
	defer C.free(unsafe.Pointer(cu))

	return togo(C.PyUnicode_FromString(cu))
}

// 数值类型 PyObject 转换 go 数值类型
func (pyObject *PyObject) Number() int {
	// a := *(*int)(unsafe.Pointer(o.ob_type))
	return int(C.PyLong_AsLong(toc(pyObject)))
}
func (pyObject *PyObject) Long() int {
	// a := *(*int)(unsafe.Pointer(o.ob_type))
	return int(C.PyLong_AsLong(toc(pyObject)))
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

// 获取 PyObject 的 __dir__()
func (pyObject *PyObject) Dir() string {
	return pyObject.CallMethod("__dir__").Repr()
}

// 获取 PyObject 的 __name__
func (pyObject *PyObject) Name() string {
	return pyObject.GetAttrString("__name__").Str()
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
// 	log.Debug().Msg(_bytes)
// 	log.Debug().Msg(C.PyBytes_Size(_bytes))
// 	if _bytes == nil {
// 		return 0
// 	} else {
// 		return int(C.PyBytes_Size(_bytes))
// 	}
// }

// Unicode 类型的 PyObject 方法
// 类型识别与数据提取
func (pyObject *PyObject) UTF8() string {
	_cutf8 := C.PyUnicode_AsUTF8(toc(pyObject))
	return C.GoString(_cutf8)
}
