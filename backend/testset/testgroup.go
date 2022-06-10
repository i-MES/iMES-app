package testset

import (
	"context"
	"fmt"
	"runtime"

	py "github.com/i-mes/imes-app/backend/python"
	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

// 测试组
type TestGroup struct {
	Id        int        `json:"id"`
	Title     string     `json:"title"`
	Desc      string     `json:"desc"`
	FileName  string     `json:"filename"`
	ClassName string     `json:"classname"`
	TestItems []TestItem `json:"testItems"`
}

func (tg *TestGroup) Run(ctx context.Context) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	wails.LogDebug(ctx, "--------- start testgroup")
	wails.LogTrace(ctx, fmt.Sprint(tg))

	if !py.Py_IsInitialized() {
		py.Py_Initialize()
	}
	_gil := py.PyGILState_Ensure()
	defer py.PyGILState_Release(_gil)
	wails.LogDebug(ctx, "Get Python GIL lock")

	_mod := py.PyImport_ImportFile(tg.FileName)
	if _mod == nil {
		panic("mod error")
	} else {
		defer _mod.DecRef()
	}

	_class := _mod.GetAttrString(tg.ClassName)
	if _class != nil {
		for _, ti := range tg.TestItems {
			_ret := _class.CallMethod(ti.FuncName)
			if _ret == nil {
				wails.LogError(ctx, fmt.Sprintf("Run TI Error: %s-%s", tg.ClassName, ti.FuncName))
			} else {
				wails.LogDebug(ctx, fmt.Sprintf("Run TI Pass: %s-%s", tg.ClassName, ti.FuncName))
			}
		}
	}
}
