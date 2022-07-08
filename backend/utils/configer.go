package utils

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

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

func getConfigFilePath(dataType string) string {
	// JSON
	if dataType == "" {
		return ""
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
		return filePath
	} else if os.IsNotExist(err) {
		// 文件不存在
		return ""
	} else {
		fmt.Println("其他 Error")
		return ""
	}
}

// 读取序列化的字符串（不做解析）
func InputConfigData(dataType string) ([]byte, error) {
	if fp := getConfigFilePath(dataType); fp != "" {
		if data, err := os.ReadFile(fp); err != nil {
			log.Fatalf("ReadFile error: %v", err)
			return nil, err
		} else {
			return data, nil
		}
	} else {
		return nil, nil
	}
}

/*
监控 Config文件 和源代码文件之间的修改时间，以触发回调。
入参：
	dataType: config 文件类型，即：config/<xyz>.json 中的 xyz
	srcs:	源代码文件路径，可从 TestClass 中获取，可从文件夹中提取
	newcallback、oldcallback: 翻转时触发一次，不是每次检查到都触发。
返回值：
    ticker：可用于结束协程
*/
func StartSyncMonitor(dataType string, srcs []string, newcallback func(string), oldcallback func()) *time.Ticker {
	ticker := time.NewTicker(time.Second * 3)
	fmt.Println("create monitor:", srcs)
	// 创建协程处理 Monitor，使用 ticker.Stop() 结束协程
	go func(ch <-chan time.Time) {
		defer ticker.Stop()
		newerfound := false
		oncecallback := false
		var newsrcfile string
		// chan + for...range 结构可以实现：
		// 1. chan 每次收到数据，执行一次 Monitor 操作。—— 可定时向 chan 喂数据
		// 2. chan 若被关闭，for 循环结束，即退出协程，否则协程一直存在。
		for t := range ch {
			if cf := getConfigFilePath(dataType); cf != "" {
				// os:
				// 		type FileInfo = fs.FileInfo
				// 		func Stat(name string) (FileInfo, error)
				// fs:
				// 		func Stat(fsys FS, name string) (FileInfo, error)
				if cfinfo, err := os.Stat(cf); err == nil {
					for _, sf := range srcs {
						if sfinfo, _err := os.Stat(sf); _err == nil {
							if cfinfo.ModTime().Before(sfinfo.ModTime()) {
								// 源代码比缓存新
								fmt.Printf("[%s(%s)] newer than [%s(%s)]\n",
									sfinfo.Name(), sfinfo.ModTime(),
									cfinfo.Name(), cfinfo.ModTime())
								newerfound = true
								newsrcfile = sf
								break
							}
						}
					}

					if newerfound {
						if !oncecallback {
							fmt.Println("newcallback Ticker:", t)
							newcallback(newsrcfile)
							oncecallback = true
						}
					} else {
						if oncecallback {
							fmt.Println("oldcallback Ticker:", t)
							oldcallback()
							oncecallback = false
						}
					}
					newerfound = false // 每轮都重新置否
				}
			}
		}
	}(ticker.C)
	return ticker
}
