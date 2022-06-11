package main

import (
	"context"
	"fmt"
	"os"

	backend "github.com/i-mes/imes-app/backend"
	py "github.com/i-mes/imes-app/backend/python"
	"github.com/siddontang/go/log"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
	api *backend.Api
}

// NewApp creates a new App application struct
func NewApp(api *backend.Api) *App {
	app := &App{}
	app.api = api
	return app
}

func (a *App) loadConfigFileCallback(data *menu.CallbackData) {
	fmt.Println(data.MenuItem.Label)
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	// wails 通过这里将上下文回吐给 user
	// user 后续用该上下文(a.ctx) 与 wails runtime 交互
	a.ctx = ctx
	a.api.Context(&ctx)
	myMenu := menu.NewMenuFromItems(
		menu.SubMenu("File", menu.NewMenuFromItems(
			menu.Text("Load Config File", keys.CmdOrCtrl("o"), a.loadConfigFileCallback),
			menu.Separator(),
			menu.Text("Quit", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
				wails.Quit(a.ctx)
			}),
		)),
	)
	wails.MenuSetApplicationMenu(a.ctx, myMenu)

	// 不是 log 到文件，而是到 stdout
	envInfo := wails.Environment(a.ctx)
	wails.LogInfo(a.ctx, envInfo.BuildType)
	if envInfo.BuildType == "dev" {
		wails.LogSetLogLevel(a.ctx, logger.DEBUG)
	} else {
		wails.LogSetLogLevel(a.ctx, logger.INFO)
	}
	wails.LogInfo(a.ctx, envInfo.Platform)
	wails.LogInfo(a.ctx, envInfo.Arch)
	wd, _ := os.Getwd()
	wails.LogInfo(a.ctx, wd)

	if !py.Py_IsInitialized() {
		py.Py_Initialize()
	}
	if !py.Py_IsInitialized() {
		panic(fmt.Errorf("python: could not initialize the python interpreter"))
	}
	log.Debug("< Initilized Python.")

	// ===== CPython 高层接口 =====
	fmt.Println("=========== Try Python High Level C API")
	py.PyRun_SimpleString("import sys")
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
print("dir():", dir()) 						
print("OS.name: ", sys.modules["os"].name)
print("Path.cwd: ", Path.cwd())
print("Path.home: ", Path.home())
`)

	// ===== CPython 低层接口 =====
	fmt.Println("=========== Try Python Low Level C API")

	// python module 的 2 步动作
	//     1. load 到 sys.modules：module 已被创建，并全局可访问
	//     2. 注入到 dir()：可以用 xxxmod.xxx 访问到该 module
	// PyImport_ImportModule -- 其实只是 Load 到 sys.modules,并没有注入到 dir() —— 这也太坑爹了，困我好几天
	// PyRun_SimpleString("import xxx") -- 才会同时生效到 sys.modules & dir()
	// PyImport_AddModule -- 不会 load 和 import，会检查 sys.modules 中是否有，有则拿出，没有创建个 empty 的。
	// 这逻辑也没谁了，太 TM 让人 emo 了。

	fmt.Println("====== Understand PyImport_ImportModule")
	_dtim := py.PyImport_ImportModule("datetime")
	// defer _dtim.DecRef() // 保留住，就不删除了
	py.PyRun_SimpleString(`print(sys.modules["datetime"])`)                         // 有，可以访问到
	py.PyRun_SimpleString(`print("Now1: ",sys.modules["datetime"].datetime.now())`) // 有效
	py.PyRun_SimpleString(`print("dir():", dir())`)                                 // 没有 datetime
	py.PyRun_SimpleString(`print("Now2: ",datetime.datetime.now())`)                // 无效, NameError: name 'datetime' is not defined

	_nowfunc := _dtim.GetAttrString("datetime").GetAttrString("now")
	defer _nowfunc.DecRef()
	fmt.Println(py.PyCallable_Check(_nowfunc)) // true

	_now := _nowfunc.CallObject(nil) // call now funcution
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

	fmt.Println("====== Understand PyImport_GetModule & PyImport_AddModule")
	_dtgm := py.PyImport_GetModule("datetime")
	_dtam := py.PyImport_AddModule("datetime")
	fmt.Println(_dtim) // 指针
	fmt.Println(_dtgm) // 指针
	fmt.Println(_dtam) // 指针，dtim == dtgm == dtam

	_mod := py.PyImport_GetModule("math")
	if _mod == nil {
		fmt.Println("Module do not imported")
	} else {
		fmt.Println(py.PyCallable_Check(_mod))       // false
		fmt.Println(_mod.HasAttrString("pi"))        // true
		fmt.Println(_mod.GetAttrString("pi").Repr()) // 3.14...
	}

	_mod = py.PyImport_AddModule("random")
	// defer _mod.DecRef() // AddModule 的是借用，不需要自己维护指针
	if _mod == nil {
		fmt.Println("Module can not added")
	} else {
		fmt.Println(py.PyCallable_Check(_mod))    // false
		fmt.Println(_mod.HasAttrString("random")) // true
		_v := _mod.CallMethod("random")
		defer _v.DecRef()
		fmt.Println(_v.Repr()) // 0.xxx....
	}

	fmt.Println("====== Understand PyImport from a file")
	// 下面不生效，原因：
	// PyImport_ImportModule(): 只会 import 已经在 sys.path 路径下的 *.py 文件
	// 文档里说的仅能使用绝对路径，不是文件的绝对路径，而是module的(绝对:x.y,相对:..x.y)
	// _gpio := py.PyImport_ImportModule("/data/kevin/workspace/kproject/imes/iMES-app/testcase/python/test_gpio.py")
	// defer _gpio.DecRef()
	// 所以封装了下面函数:
	// py.PyImport_ImportFile("./testcase/python/test_gpio.py")

	fmt.Println("====== Understand PyXXX_Check")
	fmt.Println("Is PyImport_GetModuleDict's output dict type? - ", py.PyDict_Check(py.PyImport_GetModuleDict())) // true
	// python3.7 之后已废弃
	// _tstate := C.PyGILState_GetThisThreadState()
	// C.PyEval_ReleaseThread(_tstate)

	fmt.Println("Python version: ", py.Py_GetVersion())

	// Py_Initialize() 会占用 GIL，此处不放弃，其他地方抢占不到。
	py.PyEval_SaveThread()
}

// domReady is called after the front-end dom has been loaded
// domReady 在前端Dom加载完毕后调用
func (a *App) domReady(ctx context.Context) {
	// Add your action here
	// 在这里添加你的操作
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling wails.Quit.
// Returning true will cause the application to continue,
// false will continue shutdown as normal.
// beforeClose在单击窗口关闭按钮或调用runtime.Quit即将退出应用程序时被调用.
// 返回 true 将导致应用程序继续，false 将继续正常关闭。
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	py.Py_Finalize()
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
