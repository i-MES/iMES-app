package main

import (
	"context"
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

	// 不是 log 到文件，而是到 stdout
	envInfo := wails.Environment(ctx)
	if envInfo.BuildType == "dev" {
		utils.InitLog("dev")
	} else {
		utils.InitLog("rls")
	}
	time.Sleep(time.Second * 1)
	log.Info().Msg("BuildType: " + envInfo.BuildType)
	log.Info().Msg("GOOS: " + runtime.GOOS)
	log.Info().Msg("GOARCH: " + runtime.GOARCH)
	log.Info().Msg("ProcessId: " + strconv.Itoa(utils.GetProcessIdGet()))
	if wd, e := os.Getwd(); e == nil {
		log.Info().Msg("Getwd: " + wd)
	}

	// ===== CPython 启动 =====
	if !py.Py_IsInitialized() {
		// 方式1: Unix 用 :, Windows 用 ;
		// py.Py_SetProgramName("/data/kevin/workspace/zproject/pytp/venv3.10/bin/python") // 参与生成 prefix，进而生成 sys.path
		py.Py_SetProgramName("/Users/wangkevin/workspace/kproject/mes/iMES_202204/venv3.10.4/bin/python")
		// py.Py_SetPythonHome()  // 标准库搜索路径(不可 venv)
		// py.Py_SetPath("") 			// 会清空 prefix，只留下本初的入参，所以可能缺少标准库，不如由 Py_SetProgramName() 自动生成比较稳妥
		py.Py_Initialize()

		// 方式2
		// syspath := []string{
		// 	"/home/me/.pyenv/versions/3.10.4/lib/python3.10",
		// 	"/home/me/.pyenv/versions/3.10.4/lib/python3.10/lib-dynload",
		// 	"/data/kevin/workspace/zproject/pytp/pyTP/src",
		// 	"/data/kevin/workspace/zproject/pytp/venv/lib/python3.8/site-packages"}
		// py.Py_InitializeFromConfig(syspath)

		// 方式3: 后面再 sys.path.append()，这样会带上 sys.path 的默认值
		// py.Py_Initialize()
	}
	if !py.Py_IsInitialized() {
		panic(fmt.Errorf("Could not initialize the python interpreter!"))
	}

	// 验证启动情况：
	if pyName, err := py.Py_GetProgramName(); err == nil {
		log.Debug().Str("Py_GetProgramName", pyName).Send()
	}
	if pyHome, err := py.Py_GetPythonHome(); err == nil {
		log.Debug().Str("Py_GetPythonHome", pyHome).Send()
	}
	if pyPath, err := py.Py_GetPath(); err == nil {
		// 注意 pip install -e 安装的 module 是否还存在
		log.Debug().Str("Py_GetPath", pyPath).Send()
	}

	// ===== CPython 高层接口 =====
	fmt.Println("=========== Try Python High Level C API")
	py.PyRun_SimpleString("import os")
	py.PyRun_SimpleString("import sys")
	py.PyRun_SimpleString("import threading")
	py.PyRun_SimpleString("import pytp")
	py.PyRun_SimpleString("from pathlib import Path")
	// sys.modules: 已 loaded 到内存的 module
	// 建议不要直接使用 sys.modules,而是 copy 后使用
	// dir(): 可以使用的 module、变量……
	// python 启动后默认 load 了一些 module(os,sys,time,path...70+)但没有加入到 dir() 之前是不能直接使用的
	// import 后出现在 dir()，就可以直接使用
	// sys.modules["xxx"].yyy 也可以不 import 就使用
	py.PyRun_SimpleString(`
for path in sys.path:
	print(path)
_mods = sys.modules.copy() 				
for key, value in _mods.items():	
	print(value)
print("Global dir():", dir()) 						
print("OS.name: ", sys.modules["os"].name)
print("Path.cwd: ", Path.cwd())
print("Path.home: ", Path.home())
print("sys.prefix: ", sys.prefix)
print("pytp.VERSION: ", pytp.VERSION)
sys.stdout.write('******************'+'\n')
`)
	py.PyRun_SimpleString(`import debugpy`)
	py.PyRun_SimpleString(`debugpy.listen(8899)`)
	// py.PyRun_SimpleString(`debugpy.wait_for_client()`)

	// ===== CPython 低层接口 =====
	fmt.Println("=========== Try Python Low Level C API")

	// python module 的 2 步动作
	//     1. load 到 sys.modules：module 已被创建，并全局可访问
	//     2. 注入到全局命名空间 dir()：可以用 xxxmod.xxx 访问到该 module
	// PyRun_SimpleString('import xxx') -- 完成 1、2 两个步骤
	// PyImport_ImportModule('xxx') -- 只完成 1，不做 2 —— 这也太坑爹了，困我好几天
	// 			文档中有一段：
	// 			module = PyImport_ImportModule("<modulename>");
	// 			如果模块尚未被导入（即它还不存在于 sys.modules 中），这会初始化该模块；否则它只是简单地返回 sys.modules["<modulename>"] 的值。
	// 			请注意它并不会将模块加入任何命名空间 —— 它只是确保模块被初始化并存在于 sys.modules 中。
	// 			之后你就可以通过如下方式来访问模块的属性（即模块中定义的任何名称）:
	// 			attr = PyObject_GetAttrString(module, "<attrname>");
	// 			然后可以使用 attr 自己的命名空间(即：attr.__dir__() )
	// PyImport_AddModule -- 1、2 都不做，只会检查 sys.modules 中是否有，有则拿出，没有创建个 empty 的。—— 这逻辑也没谁了，太 TM 让人 emo 了。

	fmt.Println("====== Try PyImport_ImportModule")
	_mod_dt := py.PyImport_ImportModule("datetime")
	// defer _mod_dt.DecRef() // 保留住，就不删除了
	py.PyRun_SimpleString(`print(sys.modules["datetime"])`)                         // 有，可以访问到
	py.PyRun_SimpleString(`print("Now1: ",sys.modules["datetime"].datetime.now())`) // 有效
	py.PyRun_SimpleString(`print("dir():", dir())`)                                 // 全局 dir 中没有 datetime
	py.PyRun_SimpleString(`print("Now2: ",datetime.datetime.now())`)                // 无效, NameError: name 'datetime' is not defined

	_type_dt := _mod_dt.GetAttrString("datetime")
	defer _type_dt.DecRef()
	_func_now := _type_dt.GetAttrString("now")
	// defer _func_now.DecRef()
	fmt.Println(py.PyCallable_Check(_func_now)) // true，datatime 自己的命名空间(datetime.__dir__())中有 now

	_now := _func_now.CallObject(nil) // call now funcution
	defer _now.DecRef()
	fmt.Println(py.PyNumber_Check(_now))           // false, it's type is datetime.datetime class
	fmt.Println(_now.HasAttrString("microsecond")) // true

	_attr := _now.GetAttrString("year")
	defer _attr.DecRef()
	fmt.Println(_attr)                    // *PyObject
	fmt.Println(py.PyNumber_Check(_attr)) // true
	fmt.Println(py.PyLong_Check(_attr))   // true
	fmt.Println(_attr.Number())           // 2022

	// PyObject 都会有的 3 个通用属性，相当于 repr()、dir()、type()
	fmt.Println(_now.Repr()) // datetime.datetime(2022, 6, 9, 9, 52, 33, 879300)
	fmt.Println(_now.Str())  // 2022-06-09 09:52:33.879300
	fmt.Println(_now.Type()) // <class 'datetime.datetime'>

	fmt.Println("====== Try PyImport_GetModule & PyImport_AddModule")
	_mod_dt_get := py.PyImport_GetModule("datetime")
	_mod_dt_add := py.PyImport_AddModule("datetime")
	fmt.Println(_mod_dt)     // 指针
	fmt.Println(_mod_dt_get) // 指针
	fmt.Println(_mod_dt_add) // _mod_dt == _mod_dt_get == _mod_dt_add

	_mod := py.PyImport_GetModule("math")
	if _mod == nil {
		py.PyErr_Print()
		fmt.Println("Module math can not imported")
	} else {
		fmt.Println(py.PyCallable_Check(_mod))       // false
		fmt.Println(_mod.HasAttrString("pi"))        // true
		fmt.Println(_mod.GetAttrString("pi").Repr()) // 3.14...
	}

	_mod = py.PyImport_ImportModule("random")
	// defer _mod.DecRef() // AddModule 的是借用，不需要自己维护指针
	if _mod == nil {
		py.PyErr_Print()
		fmt.Println("Module random can not added")
	} else {
		fmt.Println(py.PyCallable_Check(_mod))    // false
		fmt.Println(_mod.HasAttrString("random")) // true
		_v := _mod.CallMethod("random")
		if _v == nil {
			py.PyErr_Print()
		} else {
			if !py.Version_Check("3.10") {
				defer _v.DecRef()
			}
			fmt.Println(_v.Repr()) // 0.xxx....
		}
	}

	fmt.Println("====== Try PyImport from a file")
	// PyImport_ImportModule(): 只会 import 已经在 sys.path 路径下的 *.py 文件
	// 文档里说的仅能使用绝对路径，不是文件的绝对路径，而是module的(绝对:x.y,相对:..x.y)
	// _gpio := py.PyImport_ImportModule("/data/kevin/workspace/kproject/imes/iMES-app/testcase/python/test_gpio.py")
	// defer _gpio.DecRef()
	// 所以封装了 PyImport_ImportFile 函数:
	_mod_gpio := py.PyImport_ImportFile(".", "./testcase/python/test_gpio.py")
	if _mod_gpio == nil {
		fmt.Println("Module test_gpio can not imported from file")
	} else {
		fmt.Println(_mod_gpio.Name())
		py.PyRun_SimpleString(`print(sys.modules["testcase.python.test_gpio"])`)
	}

	fmt.Println("====== Try PyXXX_Check")
	fmt.Println("Is PyImport_GetModuleDict's output dict type? - ", py.PyDict_Check(py.PyImport_GetModuleDict())) // true

	log.Debug().Str("PythonVersion", py.Py_GetVersion()).Send()
	log.Debug().Str("PythonPath", py.PyImport_GetModule("os").GetAttrString("sys").GetAttrString("path").Repr()).Send()

	py.InitLog()

	// python3.7 之后已废弃
	// _tstate := C.PyGILState_GetThisThreadState()
	// C.PyEval_ReleaseThread(_tstate)

	// Py_Initialize() 会占用 GIL，此处不放弃，其他地方抢占不到。
	a.threadState = py.PyEval_SaveThread() // 释放 GIL，将 state 置为 null，并且返回前一个 state

	// 一些基本配置
	c := utils.GetSettingConfiger()
	if !c.IsSet("usercachepath") {
		c.Set("usercachepath", utils.GetUserCacheDefaultPath())
	}

	fmt.Println(c.AllSettings())
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
