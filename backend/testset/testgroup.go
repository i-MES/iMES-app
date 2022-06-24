package testset

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
)

// 测试组
type TestGroup struct {
	Title       string      `json:"title"`
	Desc        string      `json:"desc"`
	TestClasses []TestClass `json:"testclasses"`
}

// 解析 Python 文件，提取 TestGroup
// 对应关系：
// TestGroup -- file.py
// TestClass -- class Test_XXX
// TestItem  -- func test_xxx
func ParsePython(file string) TestGroup {
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

		cname := validClass.FindStringSubmatch(line)
		if len(cname) > 1 {
			fmt.Println("Match: ", cname[1])
			tcs = append(tcs,
				TestClass{1, cname[1], "", file, cname[1], make([]TestItem, 0)})
			continue
		}
		fname := validFunc.FindStringSubmatch(line)
		if len(fname) > 1 {
			fmt.Println("Match: ", fname[1])
			_l := len(tcs) - 1
			tcs[_l].TestItems = append(tcs[_l].TestItems,
				TestItem{fname[1], fname[1], file, fname[1], 0})
		}
	}
	return TestGroup{path.Base(file), file, tcs}
}

// 每个路径对应1个 TestGroup
func ParsePythons(filepaths []string) []TestGroup {
	tgs := make([]TestGroup, 0)
	// 默认
	for _, fp := range filepaths {
		tgs = append(tgs, ParsePython(fp))
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
func (tg *TestGroup) Run(ctx context.Context) {
	for _, tc := range tg.TestClasses {
		tc.Run(ctx)
	}
}

func (tg *TestGroup) Pause(ctx context.Context) {
}

func (tg *TestGroup) Stop(ctx context.Context) {
}
