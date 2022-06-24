package imes

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/i-mes/imes-app/backend/testset"
	"github.com/i-mes/imes-app/backend/utils"
	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

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

// 读取 Pytest 兼容的 python 文件夹，生成 TestGroup 数组
func (a *Api) LoadPythonTestGroup(selectPath bool) []testset.TestGroup {
	folderpath := ""
	if selectPath {
		// 用户选择文件夹
		folderpath = a.OpenFolder("Open TestCase Folder")
	} else {
		// 使用默认文件夹
		folderpath = utils.GetAppPath() + "/testcase/python/"
	}
	filepaths, err := a.WalkMatch(folderpath, "*.py")
	if err != nil {
		panic(err)
	}
	return testset.ParsePythons(filepaths)
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

// 开始一个测试组
func (a *Api) TestGroupStart(tg testset.TestGroup) bool {
	// do the real test
	tg.Run(*a.ctx)
	return true
}

// 开始一个测试项
func (a *Api) TestItemStart(ti testset.TestItem, tg_name string) bool {
	// do the real test
	ti.Run(tg_name)
	return true
}

var logs = make([]testset.TestItemLog, 0)

// 加载日志
func (a *Api) LoadTestItemLogs(testitemId int) []testset.TestItemLog {
	// logs = append(logs,
	// 	testset.TestItemLog{1, "PASS", time.Now().Unix()},
	// 	testset.TestItemLog{1, "NG", time.Now().Unix() + 1},
	// )
	return logs
}
