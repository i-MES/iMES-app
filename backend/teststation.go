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

func (a *Api) CreateTestStationExample() {
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
	var ts TestStation
	if data, err := utils.InputConfigData("station"); err == nil {
		_ts := []byte(json.Get(data).ToString())
		err := json.Unmarshal(_ts, &ts)
		if err != nil {
			fmt.Println(err)
		}
	}
	return ts
}
