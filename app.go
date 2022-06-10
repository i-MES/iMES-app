package main

/*
#cgo pkg-config: python3-embed
#define PY_SSIZE_T_CLEAN
#include <Python.h>
*/
import "C"
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

	// 用高层接口实现
	// PyImport_ImportModule 其实只是 LoadModule，只注入到 sys.modules,并没有注入到 dir() —— 这也太坑爹了，困我好几天
	// PyRun_SimpleString("import xxx") 才会同时生效到 sys.modules & dir()
	_dtim := py.PyImport_ImportModule("datetime")
	defer _dtim.DecRef()
	py.PyRun_SimpleString(`print("dir():", dir())`)
	py.PyRun_SimpleString(`print(sys.modules["datetime"])`)
	py.PyRun_SimpleString(`print("Now1: ",sys.modules["datetime"].datetime.now())`) // 有效
	py.PyRun_SimpleString(`print("Now2: ",datetime.datetime.now())`)                // 无效, NameError: name 'datetime' is not defined

	// 再用底层接口实现一下
	_dtam := py.PyImport_AddModule("datetime") // AddModule 我理解为 GetModule
	defer _dtam.DecRef()
	fmt.Println(_dtim)                           // 指针
	fmt.Println(_dtam)                           // 指针，且 dtim == dtam
	fmt.Println(py.PyCallable_Check(_dtam))      // false
	fmt.Println(_dtam.HasAttrString("datetime")) // true

	_nowfunc := _dtam.GetAttrString("datetime").GetAttrString("now")
	defer _nowfunc.DecRef()
	fmt.Println(py.PyCallable_Check(_nowfunc)) // true

	_now := _nowfunc.CallObject(nil) // call now funcution
	defer _now.DecRef()
	fmt.Println(py.PyNumber_Check(_now))           // false, it is datetime.datetime class
	fmt.Println(_now.HasAttrString("microsecond")) // true

	_attr := _now.GetAttrString("year")
	defer _attr.DecRef()
	fmt.Println(_attr)                    // *PyObject
	fmt.Println(py.PyNumber_Check(_attr)) // true
	fmt.Println(py.PyLong_Check(_attr))   // true
	fmt.Println(_attr.Number())           // 2022
	fmt.Println(_now.Repr())              // datetime.datetime(2022, 6, 9, 9, 52, 33, 879300)
	fmt.Println(_now.Str())               // 2022-06-09 09:52:33.879300
	fmt.Println(_now.Type())              // <class 'datetime.datetime'>

	// 下面不生效，原因：
	// PyImport_ImportModule(): 只会 import 已经在 sys.path 路径下的 *.py 文件
	// 文档里说的仅能使用绝对路径，不是文件的绝对路径，而是module的(绝对:x.y,相对:..x.y)
	// _gpio := py.PyImport_ImportModule("/data/kevin/workspace/kproject/imes/iMES-app/testcase/python/test_gpio.py")
	// defer _gpio.DecRef()
	// 所以封装了下面函数:
	// py.PyImport_ImportFile("./testcase/python/test_gpio.py")

	fmt.Println("PyImport_GetModuleDict is dict: ", py.PyDict_Check(py.PyImport_GetModuleDict())) // true
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
