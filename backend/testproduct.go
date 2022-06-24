package imes

import (
	"fmt"

	"github.com/i-mes/imes-app/backend/utils"
)

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
	_data := []byte(json.Get(utils.InputConfigData("productions")).ToString())
	err := json.Unmarshal(_data, &data)
	if err == nil {
		return data
	} else {
		fmt.Println(err)
		return nil
	}
}
