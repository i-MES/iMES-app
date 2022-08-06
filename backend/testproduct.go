package imes

import (
	"github.com/i-mes/imes-app/backend/utils"
	"github.com/rs/zerolog/log"
)

// 产品
type TestProduction struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

func (a *Api) CreateTestProductionExample() {
	// ...
}

func (a *Api) SaveTestProductions(data []TestProduction) {
	// ...
}

func (a *Api) LoadTestProductions() []TestProduction {
	if data, err := utils.InputConfigData("productions"); err == nil {
		// 找到 json 文件
		// 首先加载，然后与 读取 python 文件中的 tc 和 tg
		var tp []TestProduction
		_tp := []byte(json.Get(data).ToString())
		err := json.Unmarshal(_tp, &tp)
		if err == nil {
			return tp
		} else {
			log.Error().Stack().Err(err).Send()
			return nil
		}
	}
	return nil
}
