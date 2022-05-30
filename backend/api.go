package imes

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/i-mes/imes-app/backend/testset"
	jsoniter "github.com/json-iterator/go"
	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

var imesContext *context.Context

func ImesBind(ctx *context.Context) {
	imesContext = ctx
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

//Api struct to hold wails runtime for all Api implementations
type Api struct {
	// conf     map[interface{}]interface{}
	counter int
}

func (a *Api) InitCounter() {
	a.counter = 0
}
func (a *Api) GetCounter() int {
	return a.counter
}

func (a *Api) AddCounter(c int) int {
	a.counter += c
	return a.counter
}
func (a *Api) OpenGithub() {
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

func (a *Api) MsgDialog(msg string) {
	selection, _ := wails.MessageDialog(*imesContext, wails.MessageDialogOptions{
		Title:   "Infomation",
		Message: msg,
		Buttons: []string{"close"},
	})
	if selection == "close" {
		return
	} else {
		return
	}
}

// 通过对话框 UI 得到用户选择
func (a *Api) OpenFile(title, filePattern string) string {
	if title == "" {
		title = "Open Config File"
	}
	if filePattern == "" {
		filePattern = "*.json"
	}
	_opt := wails.OpenDialogOptions{
		DefaultDirectory: "./",
		Title:            title,
		Filters:          []wails.FileFilter{{DisplayName: "File Filter", Pattern: filePattern}},
	}
	selectedFile, err := wails.OpenFileDialog(*imesContext, _opt)
	if err != nil {
		log.Panic("Error on file opening", err.Error())
	}
	return selectedFile
}

// 通过对话框 UI 得到用户选择
func (a *Api) OpenFolder(title string) string {
	if title == "" {
		title = "Open Config Folder"
	}
	_opt := wails.OpenDialogOptions{
		DefaultDirectory: GetAppPath(),
		Title:            title,
	}
	selectedFolder, err := wails.OpenDirectoryDialog(*imesContext, _opt)
	if err != nil {
		log.Panic("Error on folder opening", err.Error())
	}
	return selectedFolder
}

func (a *Api) WalkMatch(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil

}

// 产品
type TestProduction struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

func (a *Api) InitTestProductions() {
	// ...
}
func (a *Api) SaveTestProductions(data []TestProduction) {
	// ...
}

func (a *Api) LoadTestProductions() []TestProduction {
	var data []TestProduction
	_data := []byte(json.Get(InputConfigData("productions")).ToString())
	err := json.Unmarshal(_data, &data)
	if err == nil {
		return data
	} else {
		fmt.Println(err)
		return nil
	}
}

// 工序
type TestStage struct {
	Id       int    `json:"id"`
	ProdId   int    `json:"pid"`
	Title    string `json:"title"`
	Desc     string `json:"desc"`
	Sequence int    `json:"sequence"`
}

func (a *Api) InitTestStage() {
	data := make([]TestStage, 0)
	data = append(data,
		TestStage{01, 1, "Boot测试", "测试工序", 2},
		TestStage{02, 1, "上电测试", "测试工序", 3},
		TestStage{03, 1, "烧录版本", "测试工序", 1},
		TestStage{04, 1, "校准测试", "测试工序", 1},
		TestStage{05, 1, "网络测试", "测试工序", 1},
		TestStage{11, 2, "Boot测试", "测试工序", 2},
		TestStage{12, 2, "上电测试", "测试工序", 3},
		TestStage{13, 2, "盐雾版本", "测试工序", 1},
		TestStage{14, 2, "颜色测试", "测试工序", 1},
		TestStage{15, 2, "跌落测试", "测试工序", 1},
		TestStage{21, 3, "网络测试", "测试工序", 2},
		TestStage{22, 3, "电路测试", "测试工序", 3},
		TestStage{23, 3, "盐雾版本", "测试工序", 1},
		TestStage{41, 4, "综合测试", "测试工序", 1},
	)
	a.SaveTestStages(data)
}
func (a *Api) SaveTestStages(data []TestStage) {
	_data := make(map[string]interface{})
	_data["stages"] = data
	OutputConfigData(_data)
}

func (a *Api) LoadTestStages() []TestStage {
	var data []TestStage
	_data := []byte(json.Get(InputConfigData("stages")).ToString())
	err := json.Unmarshal(_data, &data)
	if err == nil {
		return data
	} else {
		fmt.Println(err)
		return nil
	}
}

// 工位（允许支持多个测试工序）
type TestStation struct {
	Id                  int    `json:"id"`
	Title               string `json:"title"`
	Desc                string `json:"desc"`
	EnabledTestStageIds []int  `json:"enabledTestStageIds"`
	ActivedTestStageIds []int  `json:"activedTestStageIds"`
}

func (a *Api) InitTestStation() {
	a.SaveTestStation(TestStation{
		Id:                  789,
		Title:               "Station1",
		Desc:                "一个非常好用的工位",
		EnabledTestStageIds: []int{147, 369, 248},
		ActivedTestStageIds: []int{369, 248},
	})
}
func (a *Api) SaveTestStation(data TestStation) {
	_data := make(map[string]interface{})
	_data["station"] = data
	OutputConfigData(_data)
}

// 获取工位信息，通常即本机
func (a *Api) LoadTestStation() TestStation {
	var data TestStation
	_data := []byte(json.Get(InputConfigData("station")).ToString())
	err := json.Unmarshal(_data, &data)
	if err != nil {
		fmt.Println(err)
	}
	return data
}

// 被测实体
type TestEntity struct {
	Ip   []int    `json:"ip"`
	Code string   `json:"code"` // 条码
	Tags []string `json:"tags"`
}

func (a *Api) InitTestEntity() {
	a.SaveTestEntity([]TestEntity{{
		Ip:   []int{127, 0, 0, 1},
		Code: "foobar123",
		Tags: []string{"高端PC"},
	}})
}
func (a *Api) SaveTestEntity(data []TestEntity) {
	_data := make(map[string]interface{})
	_data["entity"] = data
	OutputConfigData(_data)
}
func (a *Api) LoadTestEntity() []TestEntity {
	var data []TestEntity
	_data := []byte(json.Get(InputConfigData("entity")).ToString())
	err := json.Unmarshal(_data, &data)
	if err != nil {
		fmt.Println(err)
	}
	return data
}

func (a *Api) ConnectTestEntity(ip []int) bool {
	if len(ip) == 4 {
		fmt.Println("IP V4")
	} else if len(ip) == 6 {
		fmt.Println("IP V6")
	} else {
		fmt.Println("Invallied arg")
	}
	return true
}

func (a *Api) LoadPythonTestSet() []testset.TestGroup {
	folderpath := ""
	if true {
		folderpath = GetAppPath() + "/testcase/python/"
	} else {
		// 用户选择
		folderpath = a.OpenFolder("Open TestCase Folder")
	}
	filepaths, err := a.WalkMatch(folderpath, "*.py")
	if err != nil {
		panic(err)
	}
	p := new(testset.Parser)
	tgs := make([]testset.TestGroup, 0)
	for _, fp := range filepaths {
		fmt.Println("tgs len:", len(tgs))
		tgs = append(tgs, (p.ParsePython(len(tgs), fp))...)
	}
	return tgs
}

// func (a *Api) SaveTestGroup(data []testset.TestGroup) {
// 	_data := make(map[string]interface{})
// 	_data["testgroup"] = data
// 	OutputConfigData(_data)
// }
// func (a *Api) LoadTestGroup() []testset.TestGroup {
// 	var data []testset.TestGroup
// 	_data := []byte(json.Get(InputConfigData("testgroup")).ToString())
// 	err := json.Unmarshal(_data, &data)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	return data
// }

// func (a *Api) InitTestItems() {
// 	data := make([]testset.TestItem, 0)
// 	data = append(data,
// 		testset.TestItem{"MCU Test", "MCU Test...", "test_mcu", 1},
// 		testset.TestItem{"Memory Test", "Memory Test...", "test_memory", 2},
// 		testset.TestItem{"Network Test", "Network Test...", "test_network", 3},
// 	)
// 	a.SaveTestItems(data)
// }

// func (a *Api) SaveTestItems(data []testset.TestItem) {
// 	_data := make(map[string]interface{})
// 	_data["testitem"] = data
// 	OutputConfigData(_data)
// }
// func (a *Api) LoadTestItems() []testset.TestItem {
// 	var data []testset.TestItem
// 	_data := []byte(json.Get(InputConfigData("testitem")).ToString())
// 	err := json.Unmarshal(_data, &data)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	return data
// }

// 开始一个测试项
func (a *Api) TestItemStart(id int) bool {
	// do the real test

	// add the log
	wails.EventsEmit(*imesContext, "testitemlog", testset.TestItemLog{1, "PASS", time.Now().Unix()})
	return true
}

var logs = make([]testset.TestItemLog, 0)

// 加载日志
func (a *Api) LoadTestItemLogs(testitemId int) []testset.TestItemLog {
	return append(logs,
		testset.TestItemLog{1, "PASS", time.Now().Unix()},
		testset.TestItemLog{1, "NG", time.Now().Unix() + 1},
	)
}
