package python

import (
	"errors"
	"fmt"
	"testing"

	"github.com/i-mes/imes-app/backend/utils"

	"github.com/stretchr/testify/require"
)

func TestPython(t *testing.T) {

	setup := func(t *testing.T) (func(), func()) {
		if !Py_IsInitialized() {
			Py_Initialize()
		}

		cleanup := func() {
			if Py_IsInitialized() {
				Py_Finalize()
			}
		}

		showinfo := func() {

			pyName, err := Py_GetProgramName()
			require.NoError(t, err)
			require.Contains(t, pyName, "python") // pyName 中包含 python 字符串
			t.Log("Python ProgramName: ", pyName)

			home, e := Py_GetPythonHome()
			require.NoError(t, e)
			t.Log("Python Home: ", home)

			p, ee := Py_GetPath()
			require.NoError(t, ee)
			t.Log("Python sys.path: ", p)
			t.Log("Python Version: ", Py_GetVersion())
		}
		return cleanup, showinfo
	}

	// 使用 SetProgramName 验证 python 虚拟环境的加载
	// 首先在 iMES-App 中的 setting 页面中配置虚拟环境路径
	// 并删除全局 python 中的 pytest，在虚拟环境中安装 pytest
	t.Run("Check venv with SetProgramName", func(t *testing.T) {
		// 检查是否配置了虚拟环境路径
		c := utils.GetSettingConfiger()
		if c.IsSet("pythonvenvpath") {
			dir := c.GetString("pythonvenvpath") + "/python"
			t.Log("set python venv: ", dir)
			Py_SetProgramName(dir) // 配置虚拟环境下的 python 为主程序

			cleanup, showinfo := setup(t) // Py_SetProgramName 必须在 Py_Initialize 之前操作，并且不能反复操作
			defer cleanup()

			showinfo()

			// 验证虚拟环境中是否有 pytest，全局环境中不能有
			mod := PyImport_ImportModule("debugpy")
			defer mod.DecRef()
			if mod == nil {
				// PyErr_Print()
				require.NoError(t, errors.New("Module can not imported"))
			} else {
				v := mod.GetAttrString("version_tuple")
				t.Log(v.Repr())
			}
		} else {
			t.Log("can not find python venv")
		}
	})

	// Py_SetProgramName 验证需要在第一个，其他只能在其后面
	t.Run("Initialize and Finalize", func(t *testing.T) {
		cleanup, showinfo := setup(t)
		defer cleanup()

		showinfo()

		// 验证 PyXXX_Check 函数是否有效
		require.True(t, PyDict_Check(PyImport_GetModuleDict()))
		require.False(t, PyCallable_Check(PyImport_ImportModule("os").GetAttrString("path")))  // os.path 是属性
		require.True(t, PyCallable_Check(PyImport_ImportModule("os").GetAttrString("getpid"))) // os.getpid 是方法
	})

	// 验证环境中是否包含必要的模块
	t.Run("Has necessary module", func(t *testing.T) {
		cleanup, _ := setup(t)
		defer cleanup()

		mods := []string{"pytest", "debugpy"}
		for _, mod_name := range mods {
			if PyImport_ImportModule(mod_name) == nil {
				require.NoError(t, errors.New("can not import module "+mod_name))
			}
		}
	})

	// 验证
	// 1. PyImport_ImportModule 只会加载到 sys.modules, 不会影响 dir()
	// 2. PyRun_SimpleString 即会加载到 sys.modules，也会加入 dir()
	t.Run("sys.modules and dir", func(t *testing.T) {
		cleanup, _ := setup(t)
		defer cleanup()

		_mod_dt := PyImport_ImportModule("datetime")
		if _mod_dt == nil {
			require.NoError(t, errors.New("Can not import module datetime"))
		}
		defer _mod_dt.DecRef()

		// 验证 PyImport_ImportModule() 后 sys.modules 中已经添加，但 dir() 中无法使用：
		PyRun_SimpleString(`print(sys.modules["datetime"])`)                         // 有，可以访问到
		PyRun_SimpleString(`print("Now1: ",sys.modules["datetime"].datetime.now())`) // 有效
		PyRun_SimpleString(`print("dir():", dir())`)                                 // 全局 dir 中没有 datetime
		PyRun_SimpleString(`print("Now2: ",datetime.datetime.now())`)                // 无效, NameError: name 'datetime' is not defined

		_type_dt := _mod_dt.GetAttrString("datetime") // 此时如果想使用可以取 module 的 object.GetAttrString()
		defer _type_dt.DecRef()
		_func_now := _type_dt.GetAttrString("now") // 并且可以一层层取下去
		// defer _func_now.DecRef()
		fmt.Println(PyCallable_Check(_func_now)) // true，datatime 自己的命名空间(datetime.__dir__())中有 now

		_now := _func_now.CallObject(nil) // call now funcution
		defer _now.DecRef()
		fmt.Println(PyNumber_Check(_now))              // false, it's type is datetime.datetime class
		fmt.Println(_now.HasAttrString("microsecond")) // true

		_attr := _now.GetAttrString("year")
		defer _attr.DecRef()
		fmt.Println(_attr)                 // *PyObject
		fmt.Println(PyNumber_Check(_attr)) // true
		fmt.Println(PyLong_Check(_attr))   // true
		fmt.Println(_attr.Number())        // 2022

		// PyObject 都会有的 3 个通用属性，相当于 repr()、dir()、type()
		fmt.Println(_now.Repr()) // datetime.datetime(2022, 6, 9, 9, 52, 33, 879300)
		fmt.Println(_now.Str())  // 2022-06-09 09:52:33.879300
		fmt.Println(_now.Type()) // <class 'datetime.datetime'>

		fmt.Println("====== Try PyImport_GetModule & PyImport_AddModule")
		_mod_dt_get := PyImport_GetModule("datetime")
		_mod_dt_add := PyImport_AddModule("datetime")
		fmt.Println(_mod_dt)     // 指针
		fmt.Println(_mod_dt_get) // 指针
		fmt.Println(_mod_dt_add) // _mod_dt == _mod_dt_get == _mod_dt_add

		_mod := PyImport_GetModule("math")
		if _mod == nil {
			PyErr_Print()
			fmt.Println("Module math can not imported")
		} else {
			fmt.Println(PyCallable_Check(_mod))          // false
			fmt.Println(_mod.HasAttrString("pi"))        // true
			fmt.Println(_mod.GetAttrString("pi").Repr()) // 3.14...
		}

		_mod = PyImport_ImportModule("random")
		// defer _mod.DecRef() // AddModule 的是借用，不需要自己维护指针
		if _mod == nil {
			PyErr_Print()
			fmt.Println("Module random can not added")
		} else {
			fmt.Println(PyCallable_Check(_mod))       // false
			fmt.Println(_mod.HasAttrString("random")) // true
			_v := _mod.CallMethod("random")
			if _v == nil {
				PyErr_Print()
			} else {
				if !Version_Check("3.10") {
					defer _v.DecRef()
				}
				fmt.Println(_v.Repr()) // 0.xxx....
			}
		}
	})

	t.Run("PyImport_ImportFile", func(t *testing.T) {
		_mod_gpio := PyImport_ImportFile("../../", "./testcase/python/test_gpio.py")
		if _mod_gpio == nil {
			fmt.Println("Module test_gpio can not imported from file")
		} else {
			fmt.Println(_mod_gpio.Name())
			PyRun_SimpleString(`print(sys.modules["testcase.python.test_gpio"])`)
		}
	})
}
