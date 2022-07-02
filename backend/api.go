package imes

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/i-mes/imes-app/backend/target"
	"github.com/i-mes/imes-app/backend/utils"
	jsoniter "github.com/json-iterator/go"
	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

//Api struct to hold wails runtime for all Api implementations
type Api struct {
	// conf     map[interface{}]interface{}
	ctx *context.Context
}

func (a *Api) Context(ctx *context.Context) {
	a.ctx = ctx
}

func (a *Api) OpenGithub() {
	url := "https://github.com/i-MES"
	wails.BrowserOpenURL(*a.ctx, url)
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

// 弹框显示消息
func (a *Api) MsgDialog(msg string) {
	selection, _ := wails.MessageDialog(*a.ctx, wails.MessageDialogOptions{
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
	selectedFile, err := wails.OpenFileDialog(*a.ctx, _opt)
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
		DefaultDirectory: utils.GetAppPath(),
		Title:            title,
	}
	selectedFolder, err := wails.OpenDirectoryDialog(*a.ctx, _opt)
	if err != nil {
		log.Panic("Error on folder opening", err.Error())
	}
	return selectedFolder
}

/*
加载 TestGroup
*/
func (a *Api) LoadTestGroup(selectPath bool) []target.TestGroup {
	if data, err := utils.InputConfigData("testgroup"); err == nil {
		// 找到 json 文件
		// 首先加载，然后与 读取 python 文件中的 tc 和 tg
		var tg []target.TestGroup
		_tg := []byte(json.Get(data).ToString())
		if json.Unmarshal(_tg, &tg) != nil {
			fmt.Println("can not Unmarshal json data")
		}
		// Merge 后返回
		// Merge 策略：……
		return tg
	} else if os.IsNotExist(err) {
		// json 文件不存在，从源代码中解析获取
		folderpath := ""
		if selectPath {
			// 用户选择文件夹
			folderpath = a.OpenFolder("Open TestCase Folder")
		} else {
			// 使用默认文件夹
			folderpath = utils.GetAppPath() + "/testcase/python/"
		}
		filepaths, err := utils.WalkMatch(folderpath, "*.py")
		if err != nil {
			panic(err)
		}
		return target.ParsePythons(filepaths, 1)
	}
	return nil
}

func (a *Api) SaveTestGroup(data []target.TestGroup) {
	_data := make(map[string]interface{})
	_data["testgroup"] = data
	utils.OutputConfigData(_data)
}

// wails 会以新线程的方式开启本函数；
// 所以是每个 Entity 的 每个 Group 一个测试线程。
func (a *Api) RunTestGroup(teid string, tg target.TestGroup) bool {
	// do the real test
	tg.Run(*a.ctx, teid)
	return true
}

func (a *Api) StopTestGroup(teid string, tgid string) bool {
	// do the real test
	// tg.Run(*a.ctx)
	return true
}

// wails 会以新线程的方式开启本函数；
// 所以是每个 Entity 的 每个 Group 的每个 Item 一个测试线程。
// func (a *Api) RunTestItem(teid string, tgid string, ti target.TestItem) bool {
// 	// do the real test
// 	ti.Run(tg_name)
// 	return true
// }
// func (a *Api) StopTestItem(ti target.TestItem, tg_name string) bool {
// 	// do the real test
// 	ti.Run(tg_name)
// 	return true
// }

var logs = make([]target.TestItemLog, 0)

// 加载日志
func (a *Api) LoadTestItemLogs(testitemId int) []target.TestItemLog {
	// logs = append(logs,
	// 	target.TestItemLog{1, "PASS", time.Now().Unix()},
	// 	target.TestItemLog{1, "NG", time.Now().Unix() + 1},
	// )
	return logs
}

func (a *Api) UUID() string {
	return utils.UUID()
}

// 创建 config file example，供用户修改基础和参考
func (a *Api) CreateTargetExample() {
	target.CreateTestEntityExample()
	// target.CreateTestItemExample()
}

func (a *Api) LoadTestEntity() []target.TestEntity {
	return target.LoadTestEntity()
}

// 触发 binding 相关 struct

func (a *Api) NeedStruct(tistatus target.TestItemStatus) {
	fmt.Println("Just need these struct")
}
