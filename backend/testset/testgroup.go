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
	wails.LogDebug(ctx, "--------- start testgroup "+tg.Title)
	wails.LogTrace(ctx, fmt.Sprint(tg))

	if !py.Py_IsInitialized() {
		py.Py_Initialize()
	}
	_gil := py.PyGILState_Ensure()
	defer py.PyGILState_Release(_gil)
	wails.LogDebug(ctx, "Get Python GIL lock")

	_mod := py.PyImport_ImportFile(tg.FileName)
	if _mod == nil {
		wails.LogError(ctx, "import module error")
		return
	} else {
		defer _mod.DecRef()
	}

	wails.LogDebug(ctx,
		fmt.Sprintf("Does module %s has attr %s : %t", _mod.Name(), tg.Title, _mod.HasAttrString(tg.ClassName)))
	_class := _mod.CallMethod(tg.ClassName)
	if _class != nil {
		wails.LogDebug(ctx, _class.Repr())
		wails.LogDebug(ctx, _class.Dir())
		for _, ti := range tg.TestItems {
			wails.LogDebug(ctx, "------- start testitem "+ti.FuncName)
			_ret := _class.CallMethod(ti.FuncName)
			if _ret == nil {
				py.PyErr_Occurred()
				wails.LogError(ctx, fmt.Sprintf("Run TI Error: %s-%s", tg.ClassName, ti.FuncName))
				EmitTestItemLog(ctx, false, "NG")
			} else {
				wails.LogDebug(ctx, fmt.Sprintf("Run TI Pass: %s-%s", tg.ClassName, ti.FuncName))
				EmitTestItemLog(ctx, true, "PASS")
			}
		}
	} else {
		py.PyErr_Occurred()
		wails.LogError(ctx, "--- can not get "+tg.ClassName)
	}
}
