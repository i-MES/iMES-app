package imes

import (
	"context"
	"fmt"

	"github.com/i-mes/imes-app/backend/target"
	"github.com/i-mes/imes-app/backend/utils"
	jsoniter "github.com/json-iterator/go"
	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

//Api struct to hold wails runtime for all Api implementations
type Api struct {
	// conf     map[interface{}]interface{}
	ctx context.Context
}

func (a *Api) Context(ctx context.Context) {
	a.ctx = ctx
}

func (a *Api) OpenGithub() {
	url := "https://i-mes.github.io"
	wails.BrowserOpenURL(a.ctx, url)
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
	selection, _ := wails.MessageDialog(a.ctx, wails.MessageDialogOptions{
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
func (a *Api) SelectFile(title, filePattern string) string {
	return utils.SelectFile(&a.ctx, title, filePattern)
}

// 通过对话框 UI 得到用户选择
func (a *Api) SelectFolder(title string) string {
	return utils.SelectFolder(&a.ctx, title)
}

func (a *Api) LoadTestGroup(loadFlag string, isParseFolder, isUserSelectFolderPath bool) []target.TestGroup {
	return target.LoadTestGroup(&a.ctx, loadFlag, isParseFolder, isUserSelectFolderPath)
}

func (a *Api) StopTestGroupSyncMonitor() {
	target.StopTestGroupSyncMonitor()
}

func (a *Api) SaveTestGroup(data []target.TestGroup) {
	target.SaveTestGroup(&a.ctx, data)
}

// wails 会以新线程的方式开启本函数；
// 所以是每个 Entity 的 每个 Group 一个测试线程。
func (a *Api) RunTestGroup(teid string, tg target.TestGroup) bool {
	// do the real test
	tg.Run(&a.ctx, teid)
	return true
}

func (a *Api) StopTestGroup(teid string, tgid string) bool {
	// do the real test
	// tg.Run(*a.ctx)
	return true
}

func (a *Api) UUID() string {
	return utils.UUID()
}

// 创建 config file example，供用户修改基础和参考
func (a *Api) CreateTargetExample() {
	a.CreateTestProductionExample()  // 产品
	a.CreateTestStageExample()       // 工序
	a.CreateTestStationExample()     // 工位
	target.CreateTestEntityExample() // 被测实体
	// target.CreateTestGroupExample() // TG
	// target.CreateTestClassExample() // TC
	// target.CreateTestItemExample()  // TI
}

func (a *Api) LoadTestEntity() []target.TestEntity {
	return target.LoadTestEntity()
}

// 触发 binding 相关 struct

func (a *Api) NeedStruct(tistatus target.TestItemStatus) {
	fmt.Println("Just need these struct")
}

func (a *Api) GetStringSetting(key string) string {
	return utils.GetSettingConfiger().GetString(key)
}
func (a *Api) SetStringSetting(key, value string) {
	utils.GetSettingConfiger().Set(key, value)
}
func (a *Api) GetIntSetting(key string) int {
	return utils.GetSettingConfiger().GetInt(key)
}
func (a *Api) SetIntSetting(key string, value int) {
	utils.GetSettingConfiger().Set(key, value)
}
func (a *Api) GetUserCacheDefaultPath() string {
	return utils.GetUserCacheDefaultPath()
}
func (a *Api) ReadYaml(yamldir string) interface{} {
	return utils.ReadYaml(yamldir)
}
