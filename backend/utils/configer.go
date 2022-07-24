package utils

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/fsnotify/fsnotify"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
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
								// 源代码比 config 新
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

type Configer struct {
	Name        string
	DefaultPath string       // 版本内置数据文件路径 -- 只读、不监控
	UserPath    string       // App 自动或用户手动生成数据文件路径 -- 读写、监控并自动 Reload 到 v
	v           *viper.Viper // DefaultPath、UserPath 数据分别加载到 Level6、Level3 优先级
	userv       *viper.Viper // 与 UserPath 对标的数据，实时写入 UserPath
}

func (c *Configer) Init() {
	if c.DefaultPath != "" {
		if _, err := os.Stat(c.DefaultPath); err != nil {
			os.MkdirAll(c.DefaultPath, 0750)
		}
	}
	if c.UserPath != "" {
		if _, err := os.Stat(c.UserPath); err != nil {
			os.MkdirAll(c.UserPath, 0750)
		}
	}

	dpf := c.DefaultPath + "/" + c.Name + ".yaml"
	if _, err := os.Stat(dpf); err != nil {
		os.Create(dpf)
	}
	upf := c.UserPath + "/" + c.Name + ".yaml"
	if _, err := os.Stat(upf); err != nil {
		os.Create(upf)
	}

	c.v.SetConfigName(c.Name)
	c.v.SetConfigType("yaml")
	c.userv.SetConfigName(c.Name)
	c.userv.SetConfigType("yaml")

	// DefaultPath 的 config file 加载到第 Level6 优先级的 default data 中，不监控
	if c.DefaultPath != "" {
		_v := viper.New()
		_v.SetConfigName(c.Name)
		_v.SetConfigType("yaml")
		_v.AddConfigPath(c.DefaultPath)
		if err := _v.ReadInConfig(); err == nil {
			for k, v := range _v.AllSettings() {
				c.v.SetDefault(k, v)
			}
		} else {
			fmt.Println(err)
		}
	}

	// UserPath 中的 config file 加载到第 3 优先级的 config 中，并进行监控
	if c.UserPath != "" {
		c.v.AddConfigPath(c.UserPath)
		c.userv.AddConfigPath(c.UserPath)
		c.reloadUserConfig()
		c.userv.OnConfigChange(func(e fsnotify.Event) {
			fmt.Println("Config file changed:", e.Name, e.Op)
			c.reloadUserConfig()
		})
		c.userv.WatchConfig() // 创建线程实时监控
	}
}

// 重新加载 UserPath 下 config，不关闭也不重启已有的 Watcher
func (c *Configer) reloadUserConfig() {
	if c.UserPath != "" {
		// 刷新 c.v（ReadInConfig 会覆盖同级）
		if err := c.v.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				fmt.Println("Can not find config file: imes.yaml")
			} else {
				fmt.Println("Found config file, but other error accord")
			}
		} else {
			fmt.Println("AllSettings: ", c.v.AllSettings())
		}

		// 刷新 c.userv
		if err := c.userv.ReadInConfig(); err != nil {
			fmt.Println("c.userv reload error")
		}
	}
}

// 依次尝试写入 UserPath、DefaultPath，写入一个即退出
func (c *Configer) writeToFile() {
	if c.UserPath != "" {
		c.userv.WriteConfigAs(c.UserPath + "/" + c.Name + ".yaml")
	} else {
		fmt.Println("WriteToFile Error")
	}
}

func (c *Configer) IsSet(key string) bool {
	return c.v.IsSet(key)
}
func (c *Configer) Get(key string) interface{} {
	return c.v.Get(key)
}
func (c *Configer) GetString(key string) string {
	return c.v.GetString(key)
}
func (c *Configer) GetInt(key string) int {
	return c.v.GetInt(key)
}
func (c *Configer) AllSettings() map[string]interface{} {
	return c.v.AllSettings()
}
func (c *Configer) Unmarshal(rawVal interface{}) {
	x := c.v.AllSettings()
	fmt.Println(x)
	c.v.Unmarshal(rawVal)
}
func (c *Configer) UnmarshalKey(key string, rawVal interface{}) error {
	if c.v.IsSet(key) {
		if sub := c.v.Sub(key); sub != nil {
			sub.Unmarshal(rawVal)
			return nil
		} else {
			return fmt.Errorf("UserV has key, but can not Unmarshal")
		}
	}
	return nil
}

func (c *Configer) Set(key string, value interface{}) {
	c.v.Set(key, value)
	if c.UserPath != "" {
		c.userv.Set(key, value)
		c.writeToFile()
	}
}

var appConf *Configer

// 使用 viper 全局单例做 app configer
func GetAppConfiger() *Configer {
	home, _ := Home()
	confdir := home + "/.config/iMES-app/"
	if _, err := os.Stat(confdir); err != nil {
		os.MkdirAll(confdir, 0750)
	}

	if home, err := Home(); err == nil {
		if appConf == nil {
			appConf = &Configer{
				"imes",
				GetAppPath(),
				home + "/.config/iMES-app",
				viper.GetViper(),
				viper.New(),
			}
			appConf.Init()
		}
		return appConf
	} else {
		return nil
	}
}

// 用户输入 TestCase 的路径，本函数创建对应的 cache 路径
// ~/.cache/imes-app/testcase/<uuid hash>
// 并为其创建 <configtype>.yaml 的 viper
func CreateTestcaseConfiger(folder, configtype string) *Configer {
	if id, err := Hash(folder); err == nil {
		c := &Configer{
			configtype,
			"",
			GetAppConfiger().GetString("datacachepath") + "/" + id,
			viper.New(),
			viper.New(),
		}
		c.Init()
		return c
	}
	return nil
}
