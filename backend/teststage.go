package imes

import (
	"fmt"

	"github.com/i-mes/imes-app/backend/utils"
)

// 工序
type TestStage struct {
	Id       int    `json:"id"`
	ProdId   int    `json:"pid"`
	Title    string `json:"title"`
	Desc     string `json:"desc"`
	Sequence int    `json:"sequence"`
}

func (a *Api) CreateTestStageExample() {
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
	utils.OutputConfigData(_data)
}

func (a *Api) LoadTestStages() []TestStage {
	if data, err := utils.InputConfigData("stages"); err == nil {
		var ts []TestStage
		_ts := []byte(json.Get(data).ToString())
		err := json.Unmarshal(_ts, &ts)
		if err == nil {
			return ts
		} else {
			fmt.Println(err)
			return nil
		}
	}
	return nil
}
