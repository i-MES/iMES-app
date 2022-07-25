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
	"strings"

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
	Id           string            `json:"id"`
	Title        string            `json:"title"`
	Desc         string            `json:"desc"`
	FileName     string            `json:"filename"`
	ModulePath   string            `json:"modulepath"`
	ModuleName   string            `json:"modulename"`
	ClassName    string            `json:"classname"`
	TestItems    []TestItem        `json:"testitems"`
	Parametrizes map[string]string `json:"parametrizes"`
	Fixtures     []string          `json:"fixtures"`
}

// type PythonParser struct {
// 	R                     *bufio.Reader
// 	ClassValidRegx        string
// 	FuncValidRegx         string
// 	DocValidRegx          string
// 	DocStartValidRegx     string
// 	DocEndValidRegx       string
// 	ParameterizeValidRegx string
// 	FixtureValidRegx      string
// }

// func NewPythonParser(file string) *PythonParser {
// 	var err error
// 	f, err := os.Open(file)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer f.Close()
// 	r := bufio.NewReader(f)

// 	return &PythonParser{
// 		R:                     r,
// 		ClassValidRegx:        `^class\ *(Test.*):`,
// 		FuncValidRegx:         `^\s*def (test_.*)\(`,
// 		DocValidRegx:          `^\s*"""(.*)"""`,
// 		DocStartValidRegx:     `^\s*"""(.*)`,
// 		DocEndValidRegx:       `^\s*(.*)"""`,
// 		ParameterizeValidRegx: `^\s*@pytest.mark.parametrize\("(.*)",\s*(.*)\)`,
// 		FixtureValidRegx:      `^\s*@pytest.fixture\("(.*)",\s*(.*)\)`,
// 	}
// }

// func (p *PythonParser) ParseClass() py.PyObject {
// 	valid := regexp.MustCompile(p.ClassValidRegx)
// }
// func (p *PythonParser) ParseDocstr() string {

// }

// func (p *PythonParser) ParseFixture() py.PyObject {

// }

// 逐行扫描解析 Python 文件，提取 TestClass
// 对应关系：
// TestClass -- class Test_XXX
// TestItem  -- func test_xxx
func ParsePythonOneStep(ctx *context.Context, root, file string) []TestClass {
	f, err := os.Open(file)
	if err != nil {
		wails.LogError(*ctx, "Can not open file:"+file)
		return nil
	} else {
		defer f.Close()
	}

	// 解析 Fixture 夹具 文件
	fixs := make([]string, 0)
	confile := root + "/conftest.py"
	if _, err1 := os.Stat(confile); err1 == nil {
		validFixture := regexp.MustCompile(`^\s*@pytest.fixture`)
		validFixtureFunc := regexp.MustCompile(`^\s*def (.*)\(`)
		if cf, err2 := os.Open(confile); err2 == nil {
			defer cf.Close()
			r := bufio.NewReader(cf)
			for {
				if line, err3 := r.ReadString('\n'); err3 == io.EOF {
					fmt.Println("EoF of Fixture file")
					break
				} else if err != nil {
					fmt.Printf("error reading file %s", err)
					break
				} else {
					if validFixture.Match([]byte(line)) {
						wails.LogDebug(*ctx, "Match Fixture")
						// Fixture 在下一行
						if line, err4 := r.ReadString('\n'); err4 == nil {
							_fixs := validFixtureFunc.FindStringSubmatch(line)
							fixs = append(fixs, _fixs[1])
						}
					}
				}
			}
		}
	}
	wails.LogDebug(*ctx, fmt.Sprintf("Fixture: %s", fixs))

	// 解析 pytest 文件
	r := bufio.NewReader(f)
	tcs := make([]TestClass, 0)
	validClass := regexp.MustCompile(`^class\s*([t|T]est.*):`)
	validFunc := regexp.MustCompile(`^\s*def ([t|T]est.*)\((.*)\):\s*$`)
	validDoc := regexp.MustCompile(`^\s*"""(.*)"""`)
	validDocStart := regexp.MustCompile(`^\s*"""(.*)`)
	validDocEnd := regexp.MustCompile(`^\s*(.*)"""`)
	validParameterize := regexp.MustCompile(`^\s*@pytest.mark.parametrize\("(.*)",\s*(.*)\)`)

	modpath := root
	modnamewithsuf := strings.ReplaceAll(file[len([]rune(root))+1:], "/", ".")
	suf := path.Ext(modnamewithsuf)
	modname := strings.TrimSuffix(modnamewithsuf, suf)
	funcName := ""
	args := make([]string, 0)
	huntDoc := false
	hintDoc := false
	docStr := ""
	parametrize1 := ""
	parametrize2 := ""
	paraln := 0
	ln := 0 // line number
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			wails.LogTrace(*ctx, fmt.Sprintf("EoF of %s", file))
			break
		} else if err != nil {
			fmt.Printf("error reading file %s", err)
			break
		}
		ln += 1

		if !huntDoc { // 非 DocStr 匹配状态
			// 匹配出 Parametrize
			pnames := validParameterize.FindStringSubmatch(line)
			if len(pnames) > 2 {
				wails.LogTrace(*ctx, fmt.Sprintf("Match Parametrize: %s", pnames))
				parametrize1 = pnames[1]
				parametrize2 = pnames[2]
				paraln = ln
			}

			// 匹配出 Class 名称，并创建 TestClass
			cnames := validClass.FindStringSubmatch(line)
			if len(cnames) > 1 {
				wails.LogDebug(*ctx, fmt.Sprintf("Match class: %s", cnames[1]))
				_uuid, _ := uuid.NewUUID()
				if (ln - paraln) == 1 {
					tcs = append(tcs,
						TestClass{_uuid.String(), cnames[1], "", file, modpath, modname, cnames[1], make([]TestItem, 0),
							map[string]string{parametrize1: parametrize2}, fixs})
				} else {
					tcs = append(tcs,
						TestClass{_uuid.String(), cnames[1], "", file, modpath, modname, cnames[1], make([]TestItem, 0), nil, fixs})
				}
				continue
			}

			// 匹配出当前 Class 内 Method 名称，并进入 docstr 搜索状态
			fnames := validFunc.FindStringSubmatch(line)
			if len(fnames) > 1 {
				wails.LogDebug(*ctx, fmt.Sprintf("Match function: %v", fnames))
				// 继续匹配 DocString
				funcName = fnames[1]
				args = strings.Split(strings.ReplaceAll(fnames[2], " ", ""), ",")[1:]
				huntDoc = true
			}
		} else { // hunt DocStr 状态
			docstrs := validDoc.FindStringSubmatch(line)
			if len(docstrs) > 1 {
				wails.LogTrace(*ctx, "找到单行 DocStr")
				huntDoc = false
				docStr = docstrs[1]
			} else {
				wails.LogTrace(*ctx, "寻找多行注释")
				docstrs1 := validDocStart.FindStringSubmatch(line)
				if len(docstrs1) > 1 {
					wails.LogTrace(*ctx, "找到多行注释的首行")
					docStr = docstrs1[1]
					hintDoc = true
				} else {
					if hintDoc {
						wails.LogTrace(*ctx, "继续寻找 DocStr 尾")
						docstrs2 := validDocEnd.FindStringSubmatch(line)
						if len(docstrs2) > 1 {
							wails.LogTrace(*ctx, "找到 DocStr 尾")
							docStr += docstrs2[1]
							hintDoc = false
							huntDoc = false
						} else {
							wails.LogTrace(*ctx, "未到 DocStr 尾")
							docStr += line
						}
					} else {
						// 未能在函数名的紧挨下一行找到 docstring，退出 docstr 寻找
						huntDoc = false
					}
				}
			}

			// docstr 搜索完毕，创建 TestItem
			if !huntDoc {
				wails.LogTrace(*ctx, "hunt docstr 结束，生成 TestItem"+funcName+docStr)
				_l := len(tcs) - 1
				_uuid, _ := uuid.NewUUID()
				tcs[_l].TestItems = append(tcs[_l].TestItems,
					TestItem{_uuid.String(), funcName, docStr, args})
				docStr = ""
			}
		}
	}
	return tcs
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
	wails.LogDebug(ctx, "========== start testclass "+tc.Title)
	wails.LogDebug(ctx, "go process id: "+strconv.Itoa(utils.GetProcessIdGet()))
	wails.LogDebug(ctx, "go threading id: "+strconv.Itoa(utils.GetThreadId()))
	py.LogProcessId()
	py.LogThreadId()     // 与 go threading id 相同
	py.LogInfo(tc.Title) // 其中的 threading id 与上面 2 个不同

	if !py.Py_IsInitialized() {
		py.Py_Initialize()
	}

	// 检查 conftest.py, 找到 create_entity(name) 则创建 entity
	var entObj *py.PyObject
	fixObjs := make([]*py.PyObject, 0)

	// 处理夹具文件
	if _, _err := os.Stat(tc.ModulePath + "/conftest.py"); _err == nil {
		_modct := py.PyImport_ImportFile(tc.ModulePath, tc.ModulePath+"/conftest.py")
		if _modct == nil {
			wails.LogError(ctx, "import module error")
			py.PyErr_Print()
		} else {
			defer _modct.DecRef()
			wails.LogDebug(ctx, "conftest.py:")
			wails.LogDebug(ctx, _modct.Dir())

			// 首先尝试创建全局 Entity
			entObj = _modct.CallMethodArgs("create_entity", py.PyUnicode_FromString("199.33.33.33"))
			if entObj == nil {
				py.PyErr_Print()
			}

			// 其次实例化夹具（pytest fixture）
			// if len(tc.Fixtures) > 0 {
			// 	for _, fix := range tc.Fixtures {
			// 		obj := _modct.CallMethod(fix)
			// 		py.PyErr_Print()
			// 		fmt.Println("-=-=-=-=-=-=", fix, obj)
			// 		fixObjs = append(fixObjs, obj)
			// 	}
			// }
		}
	}

	// 运行测试用例
	// 导入 py 脚本
	_mod := py.PyImport_AddPathAndImportModule(tc.ModulePath, tc.ModuleName)
	if _mod == nil {
		wails.LogError(ctx, "import module error")
		py.PyErr_Print()
		return
	} else {
		defer _mod.DecRef()
	}

	wails.LogDebug(ctx,
		fmt.Sprintf("Does module %s has Class %s : %t", _mod.Name(), tc.Title, _mod.HasAttrString(tc.ClassName)))

	// 如果当前模块（.py）中发现 create_entity，则用于实例化 entity
	wails.LogDebug(ctx,
		fmt.Sprintf("Does module %s has Func %s : %t", _mod.Name(), "create_entity", _mod.HasAttrString("create_entity")))
	if _mod.HasAttrString("create_entity") {
		entObj = _mod.CallMethod("create_entity")
		if entObj != nil {
			defer entObj.DecRef()
		} else {
			py.PyErr_Print()
		}
	}

	// Py3 C-API 使用 PyObject_CallMethod 实例化 class
	_class := _mod.CallMethod(tc.ClassName)
	if _class != nil {
		wails.LogDebug(ctx, _class.Repr())
		wails.LogDebug(ctx, _class.Dir())
		for _, ti := range tc.TestItems {
			wails.LogDebug(ctx, fmt.Sprintf("------- start testitem %s with entity(%v), fixture(%v)", ti.Title, entObj, fixObjs))
			emit("testitemstatus", ti.Id, "started")
			// 调用对象的方法，执行具体的测试项
			var _ret *py.PyObject
			if len(ti.Args) > 0 {
				// 有参
				if entObj != nil && ti.Args[0] == "entity" {
					_ret = _class.CallMethodArgs(ti.Title, entObj)
				} else {
					wails.LogError(ctx, "entity is nil")
				}
			} else {
				// 无参
				_ret = _class.CallMethodArgs(ti.Title)
			}
			if _ret == nil {
				// py.PyErr_Print()
				wails.LogDebug(ctx, fmt.Sprintf("Run TI Error: %s\t%s", tc.ClassName, ti.Title))
				emit("testitemstatus", ti.Id, "ng")
			} else {
				wails.LogDebug(ctx, fmt.Sprintf("Run TI Pass: %s\t%s", tc.ClassName, ti.Title))
				emit("testitemstatus", ti.Id, "pass")
			}
			wails.LogDebug(ctx, fmt.Sprintf("------- end testitem %s\n", ti.Title))
		}
	} else {
		py.PyErr_Print()
		wails.LogError(ctx, "--- can not get "+tc.ClassName)
	}

	if _mod.HasAttrString("teardown") {
		_mod.CallMethodArgs("teardown")
	}
	wails.LogDebug(ctx, "========== end testclass "+tc.Title)
}

// 解析 go 文件，提取 TestClass
func ParseGolang(file string) []TestClass {
	return nil
}

func (tc *TestClass) RunGolang(ctx context.Context) {
	// todo
}

func (tc *TestClass) Run(ctx *context.Context, emit func(ename, tiid, msg string)) {
	wails.LogInfo(*ctx, path.Ext(tc.FileName))
	switch path.Ext(tc.FileName) {
	case ".py":
		tc.RunPython(*ctx, emit)
	case ".go":
		tc.RunGolang(*ctx)
	default:
		wails.LogError(*ctx, "can not run testitem")
	}

	wails.EventsEmit(*ctx, "testclassfinished", tc.Id)
}
func (tc *TestClass) Pause(ctx *context.Context) {
}

func (tc *TestClass) Stop(ctx *context.Context) {
}
