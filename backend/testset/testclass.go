package testset

import (
	"context"
	"fmt"
	"strconv"

	py "github.com/i-mes/imes-app/backend/python"
	"github.com/i-mes/imes-app/backend/utils"
	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

// 测试组
type TestClass struct {
	Id        int        `json:"id"`
	Title     string     `json:"title"`
	Desc      string     `json:"desc"`
	FileName  string     `json:"filename"`
	ClassName string     `json:"classname"`
	TestItems []TestItem `json:"testitems"`
}

func (tc *TestClass) Run(ctx context.Context) {
	// 上锁 goroutine —— 似乎并不需要
	// runtime.LockOSThread()
	// defer runtime.UnlockOSThread()

	// 上锁 python 解释器线程
	_gil := py.PyGILState_Ensure()
	defer py.PyGILState_Release(_gil)
	wails.LogDebug(ctx, "Get Python GIL lock")
	// debug info
	wails.LogDebug(ctx, "--------- start testclass "+tc.Title)
	wails.LogDebug(ctx, "go process id: "+strconv.Itoa(utils.GetProcessId()))
	wails.LogDebug(ctx, "go threading id: "+strconv.Itoa(utils.GetThreadId()))
	py.LogProcessId()
	py.LogThreadId()     // 与 go threading id 相同
	py.LogInfo(tc.Title) // 其中的 threading id 与上面 2 个不同

	if !py.Py_IsInitialized() {
		py.Py_Initialize()
	}

	// 导入 py 脚本
	_mod := py.PyImport_ImportFile(tc.FileName)
	if _mod == nil {
		wails.LogError(ctx, "import module error")
		return
	} else {
		defer _mod.DecRef()
	}

	wails.LogDebug(ctx,
		fmt.Sprintf("Does module %s has attr %s : %t", _mod.Name(), tc.Title, _mod.HasAttrString(tc.ClassName)))

	// Py3 C-API 使用 PyObject_CallMethod 实例化 class
	_class := _mod.CallMethod(tc.ClassName)
	if _class != nil {
		wails.LogDebug(ctx, _class.Repr())
		wails.LogDebug(ctx, _class.Dir())
		for _, ti := range tc.TestItems {
			wails.LogDebug(ctx, "------- start testitem "+ti.FuncName)
			_ret := _class.CallMethod(ti.FuncName)
			py.LogInfo(ti.FuncName)
			if _ret == nil {
				py.PyErr_Print()
				wails.LogError(ctx, fmt.Sprintf("Run TI Error: %s-%s", tc.ClassName, ti.FuncName))
				EmitTestItemLog(ctx, false, "NG")
			} else {
				wails.LogDebug(ctx, fmt.Sprintf("Run TI Pass: %s-%s", tc.ClassName, ti.FuncName))
				EmitTestItemLog(ctx, true, "PASS")
			}
		}
	} else {
		py.PyErr_Print()
		wails.LogError(ctx, "--- can not get "+tc.ClassName)
	}
}

func (tc *TestClass) Pause(ctx context.Context) {
}

func (tc *TestClass) Stop(ctx context.Context) {
}
