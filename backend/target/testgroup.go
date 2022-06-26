package target

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"

	"github.com/google/uuid"
	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

// 测试组
type TestGroup struct {
	Id          string      `json:"id"`
	Title       string      `json:"title"`
	Desc        string      `json:"desc"`
	TestClasses []TestClass `json:"testclasses"`
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

// strategy(策略)
// 1. 所有 py 文件生成 1 个 TestGroup
// 2. 每个 py 文件生成 1 个 TestGroup
func ParsePythons(filepaths []string, strategy int) []TestGroup {
	tgs := make([]TestGroup, 0)
	if strategy == 1 {
		if _uuid, err := uuid.NewUUID(); err == nil {
			tg := TestGroup{_uuid.String(), "", "", make([]TestClass, 0)}
			for _, fp := range filepaths {
				tg.TestClasses = append(tg.TestClasses, ParsePython(fp)...)
			}
			tgs = append(tgs, tg)
		} else {
			fmt.Println("error uuid get: ", err)
		}
	} else if strategy == 2 {
		for _, fp := range filepaths {
			if _uuid, err := uuid.NewUUID(); err == nil {
				tgs = append(tgs, TestGroup{_uuid.String(), path.Base(fp), fp, ParsePython(fp)})
			} else {
				fmt.Println("error uuid get: ", err)
			}
		}
	}
	fmt.Println("tgs len:", len(tgs))
	return tgs
}

// 解析 go 文件，提取 TestGroup
func ParseGolang(file string) error {
	return nil
}

// 保存 TestGroup 信息
func (tg *TestGroup) Save(ctx context.Context) {

}

// 在本地保存信息、网络保存信息、源码文件中寻找差异，并进行 merge
func (tg *TestGroup) Merge(ctx context.Context) {

}

// 执行 TestGroup 内 TestClass 的测试，Group 内串行，Group 间并行
func (tg *TestGroup) Run(ctx context.Context, teid string) {
	for _, tc := range tg.TestClasses {
		tc.Run(ctx, teid, tg.Id)
	}
	wails.EventsEmit(ctx, "testgroupfinished", tg.Id)
}

func (tg *TestGroup) Pause(ctx context.Context) {
}

func (tg *TestGroup) Stop(ctx context.Context) {
}
