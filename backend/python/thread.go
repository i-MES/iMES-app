package python

/*
#cgo pkg-config: python3-embed
#define PY_SSIZE_T_CLEAN
#include <Python.h>
*/
import "C"

type PyThreadState C.PyThreadState
type PyGILState C.PyGILState_STATE

// In Python 3.6 and older, this function created the GIL if it didn't exist.
// func PyEval_InitThreads() {
// 	C.PyEval_InitThreads()
// }

// 保存线程状态（python 解释器会放弃 GIL 锁）
func PyEval_SaveThread() *PyThreadState {
	return (*PyThreadState)(C.PyEval_SaveThread())
}

// 恢复线程状态（python 解释器会重新掌控 GIL）
func PyEval_RestoreThread(tstate *PyThreadState) {
	C.PyEval_RestoreThread((*C.PyThreadState)(tstate))
}

// 抢占 GIL 锁
func PyGILState_Ensure() PyGILState {
	return PyGILState(C.PyGILState_Ensure())
}

// 释放 GIL 锁
func PyGILState_Release(state PyGILState) {
	C.PyGILState_Release(C.PyGILState_STATE(state))
}
