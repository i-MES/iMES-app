package imes

import (
	"fmt"

	"github.com/i-mes/imes-app/backend/utils"
)

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
	utils.OutputConfigData(_data)
}

// 获取工位信息，通常即本机
func (a *Api) LoadTestStation() TestStation {
	var data TestStation
	_data := []byte(json.Get(utils.InputConfigData("station")).ToString())
	err := json.Unmarshal(_data, &data)
	if err != nil {
		fmt.Println(err)
	}
	return data
}
