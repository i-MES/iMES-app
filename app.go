package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/i-mes/imes-app/backend/target"

	imes "github.com/i-mes/imes-app/backend"
	py "github.com/i-mes/imes-app/backend/python"
	"github.com/i-mes/imes-app/backend/utils"
	"github.com/rs/zerolog/log"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx         context.Context
	api         *imes.Api
	menu        *menu.Menu
	threadState *py.PyThreadState
}

// 制作 AppMenu
func (a *App) applicationMenu() *menu.Menu {
	a.menu = menu.NewMenu()
	FileMenu := a.menu.AddSubmenu("File")
	FileMenu.AddText("Load Config Data Only", keys.CmdOrCtrl("o"), a.loadConfigFileCallback)
	FileMenu.AddSeparator()
	FileMenu.AddText("Quit", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
		wails.Quit(a.ctx)
	})
	if runtime.GOOS == "darwin" {
		a.menu.Append(menu.EditMenu()) // on macos platform, we should append EditMenu to enable Cmd+C,Cmd+V,Cmd+Z... shortcut
	}
	ViewMenu := a.menu.AddSubmenu("View")
	ViewMenu.AddText("...", keys.CmdOrCtrl(","), a.msgDialog)
	HelpMenu := a.menu.AddSubmenu("Help")
	HelpMenu.AddText("...", keys.CmdOrCtrl("."), a.msgDialog)
	// wails.MenuSetApplicationMenu(a.ctx, a.menu)
	return a.menu
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	// wails 通过这里将上下文回吐给 user
	// user 后续用该上下文(a.ctx) 与 wails runtime 交互
	a.ctx = ctx
	a.api.Context(ctx)

	envInfo := wails.Environment(ctx) // 不是 log 到文件，而是到 stdout

	// log 到文件
	if envInfo.BuildType == "dev" {
		utils.InitLog("dev")
	} else {
		utils.InitLog("rls")
	}
	time.Sleep(time.Second * 1)

	// log golang 开发环境
	ge := log.Info().
		Str("BuildType", envInfo.BuildType).
		Str("GOOS", runtime.GOOS).
		Str("GOARCH", runtime.GOARCH).
		Str("ProcessId", strconv.Itoa(utils.GetProcessIdGet())).
		Str("LogLevel", utils.GetGlobalLevel())
	if wd, e := os.Getwd(); e == nil {
		ge.Str("Getwd", wd)
	}
	ge.Msg("Go env info")

	// ===== CPython 启动 =====
	e := log.Info()

	if !py.Py_IsInitialized() {
		// 加载 python 环境及虚拟环境的 3 种方式：
		loadpymethod := 1
		if utils.GetSettingConfiger().IsSet("app.loadpymethod") {
			loadpymethod = utils.GetSettingConfiger().GetInt("app.loadpymethod")
		}

		switch loadpymethod {
		case 1:
			// 方式1: Py_SetProgramName()、Py_SetPythonHome、Py_SetPath —— 最终证明 Py_SetProgramName 最有效
			e.Str("load.python.by", "Py_SetProgramName()")
			if utils.GetSettingConfiger().IsSet("pythonvenvpath") {
				dir := utils.GetSettingConfiger().GetString("pythonvenvpath") + "/python"
				// Py_SetProgramName 参与生成 prefix，进而生成 sys.path，
				// 不但会补全 sys.path 默认值（标准模块），还会加上 venv 库（用户模块）
				py.Py_SetProgramName(dir)
				e.Str("python.venv.path", dir)
			}
			// py.Py_SetPythonHome()  // 设置标准库搜索路径(不能用于 venv)。
			// py.Py_SetPath("") 			// Unix 用 :, Windows 用 ;但会清空 prefix，只留下手工加入的路径，缺少标准库。
			py.Py_Initialize()
		case 2:
			// 方式2: 先 Py_Initialize 带上 sys.path 的默认值（标准库），再手工 sys.path.append()
			e.Str("load.python.by", "PySys_AppendSysPath()")
			py.Py_Initialize()
			if utils.GetSettingConfiger().IsSet("pythonvenvpath") {
				dir := utils.GetSettingConfiger().GetString("pythonvenvpath") + "/python"
				py.PySys_AppendSysPath(dir)
			}
		case 3:
			// 方式3: Py_InitializeFromConfig() —— 全部都手写，不现实
			e.Str("load.python.by", "Py_InitializeFromConfig()")
			syspath := []string{
				"/home/me/.pyenv/versions/3.10.4/lib/python3.10",
				"/home/me/.pyenv/versions/3.10.4/lib/python3.10/lib-dynload",
				"/data/kevin/workspace/zproject/pytp/pyTP/src",
				"/data/kevin/workspace/zproject/pytp/venv/lib/python3.8/site-packages"}
			py.Py_InitializeFromConfig(syspath)
		}
	}
	if !py.Py_IsInitialized() {
		e.Msg("Python env info")
		log.Error().Err(errors.New("Could not initialize the python interpreter!"))
	}

	// log python 开发环境
	e.Str("Version", py.Py_GetVersion())
	if pyName, err := py.Py_GetProgramName(); err == nil {
		e.Str("ProgramName", pyName)
	}
	if pyHome, err := py.Py_GetPythonHome(); err == nil {
		e.Str("Home", pyHome)
	}
	if pyPath, err := py.Py_GetPath(); err == nil {
		e.Str("DefaultModuleSearchPath", pyPath)                                     // 相比 sys.path，只含 默认（标准）库搜索路径
		e.Str("sys.path", py.PyImport_GetModule("sys").GetAttrString("path").Repr()) // sys.path == os.sys.path
		// 		py.PyRun_SimpleString(` // 获取 sys.path 的另一种方法
		// import sys
		// for path in sys.path:
		// 	print(path)
		// `)
		e.Str("sys.module", py.PyImport_GetModule("sys").GetAttrString("modules").CallMethod("keys").Repr())
		// 		py.PyRun_SimpleString(` // 获取 sys.modules 的另一种方法
		// _mods = sys.modules.copy() // 建议不要直接使用 sys.modules,而是 copy 后使用
		// for key, value in _mods.items():
		// 	print(value)
		// `)
	}

	// 默认 import 的不多，os、sys 就是，所以可以直接 PyImport_GetModule
	e.Str("os.name", py.PyImport_GetModule("os").GetAttrString("name").Repr())
	e.Str("sys.prefix", py.PyImport_GetModule("sys").GetAttrString("prefix").Repr())
	e.Str("dir()", py.PyImport_GetModule("sys").GetAttrString("prefix").Repr())

	// pathlib 等没有默认 import，sys.modules 中找不到，所以直接 PyImport_GetModule 失败，需要先 PyImport_ImportModule
	if m := py.PyImport_GetModule("pathlib"); m != nil {
		e.Str("Path.cwd", m.GetAttrString("Path").CallMethod("cwd").Repr())
	} else {
		if mod := py.PyImport_ImportModule("pathlib"); mod != nil {
			defer mod.DecRef()
			e.Str("Path.cwd", mod.GetAttrString("Path").CallMethod("cwd").Repr())
			e.Str("Path.home", mod.GetAttrString("Path").CallMethod("home").Repr())
		}
	}
	e.Msg("Python env info")

	// dir(): 全局命名空间，列出可以全局使用的 module、变量……
	// python 启动后默认 load 了一些 module(os,sys,time,path...70+)但没有加入到 dir() 之前是不能直接使用的
	// import 后出现在 dir()，就可以直接使用
	// 也可以不 import，而用 sys.modules["xxx"].yyy 直接使用
	py.PyRun_SimpleString(`
print("Global dir():", dir())
`)

	// import 本 app 希望直接使用的（非用户测试用例所需的）、必要的 module
	py.PyRun_SimpleString("import os")
	py.PyRun_SimpleString("import sys")
	py.PyRun_SimpleString("import threading")
	py.PyRun_SimpleString("import datetime")

	// 检查是否能够成功 import 必要的 3rd module: pytest、 debugpy...
	if mod_pytest := py.PyImport_Import("pytest"); mod_pytest == nil {
		log.Error().Err(errors.New("Can not import pytest")).Send()
	} else {
		defer mod_pytest.DecRef()
		py.PyRun_SimpleString(`import pytest`)
	}

	if mod_debugpy := py.PyImport_ImportModule("debugpy"); mod_debugpy == nil {
		log.Error().Err(errors.New("Can not import debugpy")).Send()
	} else {
		defer mod_debugpy.DecRef()
		// 开启 debugpy 调试 server
		log.Info().Str("debugpy", mod_debugpy.Repr()).Send()
		py.PyRun_SimpleString(`import debugpy`)
		py.PyRun_SimpleString(`debugpy.listen(8899)`)
		// py.PyRun_SimpleString(`debugpy.wait_for_client()`)
	}

	py.InitLog()

	// python3.7 之后已废弃
	// _tstate := C.PyGILState_GetThisThreadState()
	// C.PyEval_ReleaseThread(_tstate)

	// Py_Initialize() 会占用 GIL，此处不放弃，其他地方抢占不到。
	a.threadState = py.PyEval_SaveThread() // 释放 GIL，将 state 置为 null，并且返回前一个 state

	// 一些全局默认配置
	c := utils.GetSettingConfiger()
	if !c.IsSet("usercachepath") {
		c.Set("usercachepath", utils.GetUserCacheDefaultPath())
	}

	fmt.Println("AllSettings:", c.AllSettings())
}

// domReady is called after the front-end dom has been loaded
// domReady 在前端Dom加载完毕后调用
func (a *App) domReady(ctx context.Context) {
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling wails.Quit.
// Returning true will cause the application to continue,
// false will continue shutdown as normal.
// beforeClose在单击窗口关闭按钮或调用runtime.Quit即将退出应用程序时被调用.
// 返回 true 将导致应用程序继续，false 将继续正常关闭。
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	fmt.Println("beforeClose...")

	// 放这里没啥意义，可用于不重启APP，仅重启 python
	// py.PyEval_RestoreThread(a.threadState) // 获取GIL，并将线程状态置为tstate
	// py.Py_Finalize()

	// 测试中提示用户关闭有风险，请确认。
	return false
}

// shutdown is called at application termination
// 在应用程序终止时被调用
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
	// 在此处做一些资源释放的操作
}

type SysInfo struct {
	BuildType string `json:"buildtype"`
	Platform  string `json:"platform"`
	Arch      string `json:"arch"`
}

func (a *App) SysInfo() SysInfo {
	envInfo := wails.Environment(a.ctx)
	return SysInfo{
		BuildType: envInfo.BuildType,
		Platform:  envInfo.Platform,
		Arch:      envInfo.Arch,
	}
}

func (a *App) loadConfigFileCallback(data *menu.CallbackData) {
	target.LoadTestGroup(&a.ctx, "src", true, true)
}

func (a *App) msgDialog(data *menu.CallbackData) {
	a.api.MsgDialog(data.MenuItem.Label)
}
