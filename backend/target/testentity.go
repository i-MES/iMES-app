package target

import (
	"fmt"

	"github.com/i-mes/imes-app/backend/utils"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// 被测实体
type TestEntity struct {
	Id   string   `json:"id"`
	Ip   []int    `json:"ip"`
	Code string   `json:"code"` // 条码
	Tags []string `json:"tags"`
}

func CreateTestEntityExample() {
	SaveTestEntity([]TestEntity{{
		Id:   utils.UUID(),
		Ip:   []int{127, 0, 0, 1},
		Code: "foobar123",
		Tags: []string{"高端PC"},
	}})
}
func SaveTestEntity(data []TestEntity) {
	_data := make(map[string]interface{})
	_data["entity"] = data
	utils.OutputConfigData(_data)
}
func LoadTestEntity() []TestEntity {
	var data []TestEntity
	_data := []byte(json.Get(utils.InputConfigData("entity")).ToString())
	err := json.Unmarshal(_data, &data)
	if err != nil {
		fmt.Println(err)
	}
	return data
}

func ConnectTestEntity(ip []int) bool {
	if len(ip) == 4 {
		fmt.Println("IP V4")
	} else if len(ip) == 6 {
		fmt.Println("IP V6")
	} else {
		fmt.Println("Invallied arg")
	}
	return true
}
