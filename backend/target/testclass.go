package target

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"runtime"
	"strconv"

	"github.com/google/uuid"
	py "github.com/i-mes/imes-app/backend/python"
	"github.com/i-mes/imes-app/backend/utils"
	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

/*
TestClass
TestClass 与 TestItem 之间的关系根据代码而固化，用户不能调整。
只能开发者通过修改代码才能修改 TC 与 TI 之间的关系。
获取途径：
	1. Parse: 从 python、go 等文件中解析
	2. 没有了！
合并（Merge）：TestGroup 才需要 Merge，TestClass 和 TestItem 不需要。
*/
type TestClass struct {
	Id        string     `json:"id"`
	Title     string     `json:"title"`
	Desc      string     `json:"desc"`
	FileName  string     `json:"filename"`
	ClassName string     `json:"classname"`
	TestItems []TestItem `json:"testitems"`
}

// 解析 Python 文件，提取 TestClass
// 对应关系：
// TestClass -- class Test_XXX
// TestItem  -- func test_xxx
func ParsePython(file string) []TestClass {
	var err error
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	validClass := regexp.MustCompile(`^class\ *(.*):`)
	validFunc := regexp.MustCompile(`^\s*def (test_.*)\(`)
	tcs := make([]TestClass, 0)
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			fmt.Println("EoF of", file)
			break
		} else if err != nil {
			fmt.Printf("error reading file %s", err)
			break
		}

		// 匹配出 Class 名称
		cname := validClass.FindStringSubmatch(line)
		if len(cname) > 1 {
			fmt.Println("Match class: ", cname[1])
			_uuid, _ := uuid.NewUUID()
			tcs = append(tcs,
				TestClass{_uuid.String(), cname[1], "", file, cname[1], make([]TestItem, 0)})
			continue
		}

		// 匹配出 Function 名称
		fname := validFunc.FindStringSubmatch(line)
		if len(fname) > 1 {
			fmt.Println("Match function: ", fname[1])
			_l := len(tcs) - 1
			_uuid, _ := uuid.NewUUID()
			tcs[_l].TestItems = append(tcs[_l].TestItems,
				TestItem{_uuid.String(), fname[1], fname[1], file, fname[1], 0})
		}
	}
	return tcs
}

// 解析 go 文件，提取 TestClass
func ParseGolang(file string) []TestClass {
	return nil
}

func (tc *TestClass) RunPython(ctx context.Context, emit func(ename, tiid, msg string)) {
	// go 维度线程上锁
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// py 维度进程上锁
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
			emit("testitemstatus", ti.Id, "started")
			// 调用对象的方法，执行具体的测试项
			_ret := _class.CallMethod(ti.FuncName)
			py.LogInfo(ti.FuncName)
			if _ret == nil {
				py.PyErr_Print()
				wails.LogError(ctx, fmt.Sprintf("Run TI Error: %s\t%s", tc.ClassName, ti.FuncName))
				emit("testitemstatus", ti.Id, "ng")
			} else {
				wails.LogDebug(ctx, fmt.Sprintf("Run TI Pass: %s\t%s", tc.ClassName, ti.FuncName))
				emit("testitemstatus", ti.Id, "pass")
			}
		}
	} else {
		py.PyErr_Print()
		wails.LogError(ctx, "--- can not get "+tc.ClassName)
	}
}

func (tc *TestClass) RunGo(ctx context.Context) {
	// todo
}

func (tc *TestClass) Run(ctx context.Context, emit func(ename, tiid, msg string)) {
	wails.LogInfo(ctx, path.Ext(tc.FileName))
	switch path.Ext(tc.FileName) {
	case ".py":
		tc.RunPython(ctx, emit)
	case ".go":
		tc.RunGo(ctx)
	default:
		wails.LogError(ctx, "can not run testitem")
	}

	wails.EventsEmit(ctx, "testclassfinished", tc.Id)
}
func (tc *TestClass) Pause(ctx context.Context) {
}

func (tc *TestClass) Stop(ctx context.Context) {
}
