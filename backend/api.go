package imes

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	jsoniter "github.com/json-iterator/go"
	wails "github.com/wailsapp/wails/v2/pkg/runtime"
	"gopkg.in/yaml.v3"
)

var imesContext *context.Context

func ImesBind(ctx *context.Context) {
	imesContext = ctx
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

//Api struct to hold wails runtime for all Api implementations
type Api struct {
	conf    map[interface{}]interface{}
	confstr []byte
}

func (s *Api) OpenFile(Hash string) bool {
	fmt.Println("OpenFile")
	return true
}

func (s *Api) OpenLink(Link string) bool {
	fmt.Println("OpenLink")
	return true
}

func (s *Api) OpenLog() bool {
	fmt.Println("OpenLog")
	return true
}

func (s *Api) OpenFolder(Hash string) bool {
	fmt.Println("OpenFolder")
	return true
}

func (s *Api) OpenGithub() {
	url := "https://github.com/i-MES"
	wails.BrowserOpenURL(*imesContext, url)
	// var err error
	// switch runtime.GOOS {
	// case "linux":
	// 	err = exec.Command("xdg-open", url).Start()
	// case "windows":
	// 	err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	// case "darwin":
	// 	err = exec.Command("open", url).Start()
	// default:
	// 	err = fmt.Errorf("unsupported platform")
	// }
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

// 产品
type TestProduction struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

// 工序
type TestStep struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Desc     string `json:"desc"`
	Sequence int    `json:"sequence"`
}

var tss = make([]TestStep, 0)

// 加载测试工序
func (s *Api) LoadTestSteps() []TestStep {
	if len(tss) == 0 {
		tss = append(tss,
			TestStep{369, "工序2", "测试工序2", 2},
			TestStep{248, "工序3", "测试工序3", 3},
			TestStep{147, "工序1", "测试工序1", 1},
		)
	}
	return tss
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
func (s *Api) GetTestStationInfo() TestStation {
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

func (s *Api) ConnectTestEntity(ip []int) bool {
	if len(ip) == 4 {
		fmt.Println("IP V4")
	} else if len(ip) == 6 {
		fmt.Println("IP V6")
	} else {
		fmt.Println("Invallied arg")
	}
	return true
}

func (s *Api) GetActivedTestEntity() []TestEntity {
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

func (s *Api) LoadTestGroup(stepId int, stationId int, entityId int) []TestGroup {
	if len(tgs) == 0 {
		tgs = append(tgs,
			TestGroup{1, "Group1", "测试组1", []int{1, 2}},
			TestGroup{1, "Group2", "测试组2", []int{2, 3}},
		)
	}
	return tgs
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

// Load testitems from a file
func (s *Api) LoadTestItems(path string) []TestItem {
	if len(tis) == 0 {
		tis = append(tis,
			TestItem{1, "MCU Test", "MCU Test...", "test_mcu", 1},
			TestItem{2, "Memory Test", "Memory Test...", "test_memory", 2},
			TestItem{3, "Network Test", "Network Test...", "test_network", 3},
		)
	}

	return tis
}

// 开始一个测试项
func (s *Api) TestItemStart(id int) bool {
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

var logs = make([]TestItemLog, 0)

// 加载日志
func (s *Api) LoadTestItemLogs(testitemId int) []TestItemLog {
	return append(logs,
		TestItemLog{1, "PASS", time.Now().Unix()},
		TestItemLog{1, "NG", time.Now().Unix() + 1},
	)
}

var counter = 0

func (s *Api) AddCounter() int {
	counter += 1
	fmt.Println(counter)
	return counter
}

func (s *Api) LoadCounter() int {
	return counter
}

func (a *Api) OpenConfigFile() string {
	_opt := wails.OpenDialogOptions{
		DefaultDirectory: "./",
		Title:            "Open Config File",
		Filters:          []wails.FileFilter{{DisplayName: "Config File", Pattern: "*.json"}},
	}
	selectedFile, err := wails.OpenFileDialog(*imesContext, _opt)
	if err != nil {
		log.Panic("Error on file opening", err.Error())
	}
	return selectedFile
}

func (a *Api) OpenConfigFolder() string {
	_opt := wails.OpenDialogOptions{
		DefaultDirectory: "./config/",
		Title:            "Open Config Folder",
	}
	selectedFolder, err := wails.OpenDirectoryDialog(*imesContext, _opt)
	if err != nil {
		log.Panic("Error on folder opening", err.Error())
	}
	return selectedFolder
}

func (a *Api) LoadYamlConfigData(filePath string) bool {
	m := make(map[interface{}]interface{})
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("ReadFile error: %v", err)
		return false
	}
	err = yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		log.Fatalf("yaml.Unmarshal error: %v", err)
		return false
	}
	a.conf = m
	fmt.Printf("--- conf:\n%v\n\n", a.conf)
	return true
}

func (a *Api) LoadJsonConfigData(filePath string) bool {
	if filePath == "" {
		return false
	}
	matched, err := regexp.Match(`/.*`, []byte(filePath))
	if !matched {
		log.Fatalf("config file path(%v) invalled, err: %v", filePath, err)
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("ReadFile error: %v", err)
		return false
	} else {
		a.confstr = data
		// fmt.Printf("--- conf:\n%v\n\n", a.confstr)
	}
	return true
}

// func (a *Api) GetYamlConfig(keys ...interface{}) jsoniter.Any {
// 	if v, ok := a.conf[keys[0]]; ok {
// 		if len(keys) == 1 {
// 			fmt.Printf("find config %v(type:%v): %v\n", keys[0], reflect.TypeOf(v), v)
// 			return v, true
// 		} else {
// 			return a.GetYamlConfig(keys[1:])
// 		}
// 	}
// 	return "", false
// }

func (a *Api) GetYamlProductions() []TestProduction {
	var tps []TestProduction
	_tps := []byte(json.Get(a.confstr, "productions").ToString())
	err := json.Unmarshal(_tps, &tps)
	if err != nil {
		fmt.Println(tps)
		// tps := tps.([]interface{})
		// var tp TestProduction
		// for i := 0; i < len(tps); i++ {
		// 	_tp, err := yaml.Marshal(tps[i])
		// 	if err == nil {
		// 		fmt.Println(_tp)
		// 		err := yaml.Unmarshal(_tp, &tp)
		// 		if err == nil {
		// 			fmt.Println(tp)
		// 		} else {
		// 			fmt.Println("oops")
		// 		}
		// 	}
		// }
	}
	return nil
}

func (a *Api) GetJsonProductions() []TestProduction {
	var tps []TestProduction
	_tps := []byte(json.Get(a.confstr, "productions").ToString())
	err := json.Unmarshal(_tps, &tps)
	if err == nil {
		// fmt.Println(reflect.TypeOf(tps), tps)
		return tps
	} else {
		fmt.Println(err)
	}
	return nil
}
