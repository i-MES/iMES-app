package imes

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"time"

	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

var imesContext *context.Context

func ImesBind(ctx *context.Context) {
	imesContext = ctx
}

//Middleware struct to hold wails runtime for all middleware implementations
type Middleware struct {
}

func (s *Middleware) OpenFile(Hash string) bool {
	fmt.Println("OpenFile")
	return true
}

func (s *Middleware) OpenLink(Link string) bool {
	fmt.Println("OpenLink")
	return true
}

func (s *Middleware) OpenLog() bool {
	fmt.Println("OpenLog")
	return true
}

func (s *Middleware) OpenFolder(Hash string) bool {
	fmt.Println("OpenFolder")
	return true
}

func (s *Middleware) OpenGithub() {
	var err error
	url := "https://github.com/i-MES"
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

// 工序
type TestStep struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Desc     string `json:"desc"`
	Sequence int    `json:"sequence"`
}

// 加载测试工序
func (s *Middleware) LoadTestSteps() []TestStep {
	tps := make([]TestStep, 0)
	tps = append(tps,
		TestStep{147, "工序1", "测试工序1", 1},
		TestStep{369, "工序2", "测试工序2", 2},
		TestStep{248, "工序3", "测试工序3", 3},
	)
	return tps
}

// 工位（允许支持多个测试工序）
type TestStation struct {
	Id                 int    `json:"id"`
	Title              string `json:"title"`
	Desc               string `json:"desc"`
	EanbledTestStepIds []int  `json:"enabledTestStepIds"`
	ActivedTestStepIds []int  `json:"activedTestStepIds"`
}

// 获取工位信息，通常即本机
func (s *Middleware) GetTestStationInfo() TestStation {
	return TestStation{
		789,
		"Station1",
		"一个非常好用的工位",
		[]int{147, 369, 248},
		[]int{369, 248},
	}
}

// 被测实体
type TestEntity struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

var tes = make([]TestEntity, 0)

func (s *Middleware) ConnectTestEntity(ip []int) bool {
	if len(ip) == 4 {
		fmt.Println("IP V4")
	} else if len(ip) == 6 {
		fmt.Println("IP V6")
	} else {
		fmt.Println("Invallied arg")
	}
	return true
}

func (s *Middleware) GetActivedTestEntity() []TestEntity {
	return append(tes,
		TestEntity{1, "Entity1", "PC"},
		TestEntity{2, "Entity2", "MBP"},
		TestEntity{3, "Entity3", "OPPO"},
	)
}

// 测试组
type TestGroup struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Desc        string `json:"desc"`
	TestItemIds []int  `json:"testItemIds"`
}

var tgs = make([]TestGroup, 0)

func (s *Middleware) LoadTestGroup(stepId int, stationId int, entityId int) []TestGroup {
	return append(tgs,
		TestGroup{1, "Group1", "测试组1", []int{1, 2}},
		TestGroup{1, "Group2", "测试组2", []int{2, 3}},
	)
}

// 测试项
type TestItem struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Desc     string `json:"desc"`
	Funcname string `json:"funcname"`
	Sequence int    `json:"sequence"`
}

var tis = make([]TestItem, 0)
var tised bool = false

// Load testitems from a file
func (s *Middleware) LoadTestItems(path string) []TestItem {
	if !tised {
		tis = append(tis,
			TestItem{1, "MCU Test", "MCU Test...", "test_mcu", 1},
			TestItem{2, "Memory Test", "Memory Test...", "test_memory", 2},
			TestItem{3, "Network Test", "Network Test...", "test_network", 3},
		)
		tised = true
	}

	return tis
}

// 开始一个测试项
func (s *Middleware) TestItemStart(id int) bool {
	// do the real test

	// add the log
	wails.EventsEmit(*imesContext, "testitemlog", TestItemLog{1, "PASS", time.Now().Unix()})
	return true
}

// 测试项日志
type TestItemLog struct {
	TestItemId int    `json:"testItemId"`
	Message    string `json:"message"`
	TimeStamp  int64  `json:"timestamp"`
}

// 加载日志
func (s *Middleware) LoadTestItemLogs(testitemId int) []TestItemLog {
	logs := make([]TestItemLog, 0)
	return append(logs,
		TestItemLog{1, "PASS", time.Now().Unix()},
		TestItemLog{1, "NG", time.Now().Unix() + 1},
	)
}

var counter = 0

func (s *Middleware) AddCounter() int {
	counter += 1
	fmt.Println(counter)
	return counter
}

func (s *Middleware) LoadCounter() int {
	return counter
}
