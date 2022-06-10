package testset

import (
	"fmt"
	"runtime"

	py "github.com/i-mes/imes-app/backend/python"
)

// 测试项
type TestItem struct {
	Title    string `json:"title"`
	Desc     string `json:"desc"`
	FileName string `json:"filename"`
	FuncName string `json:"funcname"`
	Sequence int    `json:"sequence"`
}

func (ti *TestItem) Run(tg_name string) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	fmt.Println(ti)

	if !py.Py_IsInitialized() {
		py.Py_Initialize()
	}
	fmt.Println("-=")
	_gil := py.PyGILState_Ensure()
	defer py.PyGILState_Release(_gil)

	fmt.Println("-=-=")
	_mod := py.PyImport_ImportFile(ti.FileName)
	defer _mod.DecRef()
	if _mod == nil {
		panic("mod error")
	}

	pFunc := _mod.GetAttrString(ti.FuncName)
	if pFunc != nil {
		pValue := pFunc.CallObject(nil)
		if pValue != nil {
			fmt.Println(pValue)
		}
	}
	fmt.Println("Run TI: ", ti.Title, ti.Desc, ti.FuncName, ti.Sequence)
}
