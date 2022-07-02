package utils

import (
	"fmt"
	"log"
	"os"
	"regexp"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// data: [key: contents]，取 key 作为 filename
func OutputConfigData(data map[string]interface{}) {
	datatype := ""
	for k := range data {
		datatype = k
		// switch data[k].(type) {
		// case int:
		// 	fmt.Println("int")
		// case float64:
		// 	fmt.Println("float64")
		// case string:
		// 	fmt.Println("string")
		// default:
		// 	fmt.Println("default")
		// }
	}
	fmt.Println("Output config data, type: ", datatype)
	_data, _ := json.Marshal(data[datatype])
	filePath := GetAppPath() + "/config/" + datatype + ".json"
	err := os.WriteFile(filePath, _data, 0644)
	if err != nil {
		panic(err)
	}
}

// 读取序列化的字符串（不做解析）
func InputConfigData(dataType string) ([]byte, error) {
	// JSON
	if dataType == "" {
		return nil, nil
	}

	// 组合路径，并判断是否合法
	filePath := GetAppPath() + "/config/" + dataType + ".json"
	matched, err := regexp.Match(`/.*`, []byte(filePath))
	if !matched {
		log.Fatalf("config file path(%v) invalled, err: %v", filePath, err)
	}

	// 判断是否存在
	if _, err := os.Stat(filePath); err == nil {
		// 文件存在
		fmt.Println("Input config data from: ", filePath)
		data, err2 := os.ReadFile(filePath)
		if err2 != nil {
			log.Fatalf("ReadFile error: %v", err)
			return nil, err2
		} else {
			return data, nil
		}
	} else if os.IsNotExist(err) {
		// 文件不存在
		return nil, err
	} else {
		fmt.Println("其他 Error")
		return nil, nil
	}

	// YAML
	// m := make(map[interface{}]interface{})
	// data, err := os.ReadFile(filePath)
	// if err != nil {
	// 	log.Fatalf("ReadFile error: %v", err)
	// 	return false
	// }
	// err = yaml.Unmarshal([]byte(data), &m)
	// if err != nil {
	// 	log.Fatalf("yaml.Unmarshal error: %v", err)
	// 	return false
	// }
	// a.conf = m
	// fmt.Printf("--- conf:\n%v\n\n", a.conf)
	// return true
}
