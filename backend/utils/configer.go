package utils

import (
	"fmt"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
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
		// 	log.Debug().Msg("int")
		// case float64:
		// 	log.Debug().Msg("float64")
		// case string:
		// 	log.Debug().Msg("string")
		// default:
		// 	log.Debug().Msg("default")
		// }
	}
	log.Debug().Msgf("Output config data, type: %v", datatype)
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
		log.Fatal().Msgf("config file path(%v) invalled, err: %v", filePath, err)
	}

	// 判断是否存在
	if _, err := os.Stat(filePath); err == nil {
		// 文件存在
		return filePath
	} else if os.IsNotExist(err) {
		// 文件不存在
		return ""
	} else {
		log.Debug().Msg("其他 Error")
		return ""
	}
}

// 读取序列化的字符串（不做解析）
func InputConfigData(dataType string) ([]byte, error) {
	if fp := getConfigFilePath(dataType); fp != "" {
		if data, err := os.ReadFile(fp); err != nil {
			log.Fatal().Msgf("ReadFile error: %v", err)
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
	log.Debug().Msgf("create monitor: %v", srcs)
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
								log.Debug().Msgf("[%s(%s)] newer than [%s(%s)]\n",
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
							log.Debug().Msgf("newcallback Ticker: %v", t)
							newcallback(newsrcfile)
							oncecallback = true
						}
					} else {
						if oncecallback {
							log.Debug().Msgf("oldcallback Ticker: %v", t)
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
	Name          string       // yaml 文件名，后缀(.yaml)省略
	ReadOnlyPath  string       // 版本内置数据文件路径 -- 只读、不监控
	ReadWritePath string       // App 自动或用户手动生成数据文件路径 -- 读写、监控并自动 Reload 到 v
	alldata       *viper.Viper // ReadWritePath、 ReadWritePath 数据分别加载到 Level6、Level3 优先级
	rwdata        *viper.Viper // 与 ReadWritePath 对标的数据，实时写入 ReadWritePath
	wtimer        *time.Timer  // 写文件缓冲 timer
	mutex         *sync.Mutex  // 并发锁
}

func (c *Configer) Init() {
	// 创建 DefaultPath、UserPath 路径下相应文件
	if c.ReadOnlyPath != "" {
		if _, err := os.Stat(c.ReadOnlyPath); err != nil {
			os.MkdirAll(c.ReadOnlyPath, 0750)
		}
	}
	if c.ReadWritePath != "" {
		if _, err := os.Stat(c.ReadWritePath); err != nil {
			os.MkdirAll(c.ReadWritePath, 0750)
		}
	}

	dpf := c.ReadOnlyPath + "/" + c.Name + ".yaml"
	if _, err := os.Stat(dpf); err != nil {
		os.Create(dpf)
	}
	upf := c.ReadWritePath + "/" + c.Name + ".yaml"
	if _, err := os.Stat(upf); err != nil {
		os.Create(upf)
	}

	c.alldata.SetConfigName(c.Name)
	c.alldata.SetConfigType("yaml")
	c.rwdata.SetConfigName(c.Name)
	c.rwdata.SetConfigType("yaml")

	// ReadOnlyPath 的 config file 加载到第 Level6 优先级的 default data 中，不监控
	if c.ReadOnlyPath != "" {
		_v := viper.New()
		_v.SetConfigName(c.Name)
		_v.SetConfigType("yaml")
		_v.AddConfigPath(c.ReadOnlyPath)
		if err := _v.ReadInConfig(); err == nil {
			for k, v := range _v.AllSettings() {
				// 遍历并写入 viper 的 default data 中
				// c.mutex.Lock()
				c.alldata.SetDefault(k, v)
				// c.mutex.Unlock()
			}
		} else {
			log.Error().Stack().Err(err).Send()
		}
	}

	// ReadWritePath 中的 config file 加载到第 3 优先级的 config 中，并进行监控
	if c.ReadWritePath != "" {
		c.alldata.AddConfigPath(c.ReadWritePath)
		c.rwdata.AddConfigPath(c.ReadWritePath)
		c.reloadUserConfig()
		c.rwdata.OnConfigChange(func(e fsnotify.Event) {
			log.Debug().Msgf("Config file changed: %s, %v", e.Name, e.Op)
			// c.reloadUserConfig()
		})
		c.rwdata.WatchConfig() // 创建线程实时监控
	}

	// 创建一个写文件缓冲协程
	c.wtimer = time.NewTimer(time.Second * 3)
	go func(ch <-chan time.Time) {
		defer c.wtimer.Stop()
		for t := range ch {
			log.Debug().Msg("write config to file" + t.GoString())
			c.writeToFile()
		}
	}(c.wtimer.C)

	c.mutex = &sync.Mutex{}
}

// 重新加载 ReadWritePath 下 config，不关闭也不重启已有的 Watcher
func (c *Configer) reloadUserConfig() {
	if c.ReadWritePath != "" {
		// 刷新 c.v（ReadInConfig 会覆盖同级）
		if err := c.alldata.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				log.Debug().Msg("Can not find config file: imes.yaml")
			} else {
				log.Debug().Msg("Found config file, but other error accord")
			}
		}

		// 刷新 c.userv
		if err := c.rwdata.ReadInConfig(); err != nil {
			log.Debug().Msg("c.userv reload error")
		}
	}
}

// 依次尝试写入 UserPath、DefaultPath，写入一个即退出
func (c *Configer) writeToFile() {
	if c.ReadWritePath != "" {
		c.rwdata.WriteConfigAs(c.ReadWritePath + "/" + c.Name + ".yaml")
	} else {
		log.Debug().Msg("WriteToFile Error")
	}
}

func (c *Configer) Set(key string, value interface{}) {
	c.mutex.Lock()
	c.alldata.Set(key, value)
	if c.ReadWritePath != "" {
		c.rwdata.Set(key, value)
		// c.writeToFile()
		c.wtimer.Reset(time.Second * 3)
	}
	c.mutex.Unlock()
}

func (c *Configer) IsSet(key string) bool {
	return c.alldata.IsSet(key)
}
func (c *Configer) Get(key string) interface{} {
	return c.alldata.Get(key)
}
func (c *Configer) GetString(key string) string {
	return c.alldata.GetString(key)
}
func (c *Configer) GetInt(key string) int {
	return c.alldata.GetInt(key)
}
func (c *Configer) AllSettings() map[string]interface{} {
	return c.alldata.AllSettings()
}
func (c *Configer) Unmarshal(rawVal interface{}) {
	x := c.alldata.AllSettings()
	log.Debug().Msgf("allsettings: %v", x)
	c.alldata.Unmarshal(rawVal)
}
func (c *Configer) UnmarshalKey(key string, rawVal interface{}) error {
	if c.alldata.IsSet(key) {
		if sub := c.alldata.Sub(key); sub != nil {
			sub.Unmarshal(rawVal)
			return nil
		} else {
			return fmt.Errorf("UserV has key, but can not Unmarshal")
		}
	}
	return nil
}

/*
config 数据:

根据用途分 2 种

1. settings: 控制程序运行的配置数据，清空对 app 运行会有一定、甚至深刻影响
2. cache: 缓存类，用户可以随时清空，对 app 运行没有影响

根据用户分 2 种

1. app data：iMES-app 使用的数据
2. user data: 用户生成、使用的数据

这是一个二维问题域：
|          | app data                  | user data                                              |
| -------- | ------------------------- | ------------------------------------------------------ |
| Settings | 随 app 自带，用户不能修改 | 用户修改的 Setting，主要包括 app 中 setting 页面中内容 |
| Cache    | APP 运行态生成的缓存数据  | 用户生成、使用的数据，如：testgroup...                 |
*/

var settingConf *Configer

// 使用 viper 全局单例做 app configer
func GetSettingConfiger() *Configer {
	if settingConf == nil {
		if home, err := Home(); err == nil {
			dir := home + "/.config/imes-app/"
			if _, err := os.Stat(dir); err != nil {
				os.MkdirAll(dir, 0750)
			}
			settingConf = &Configer{
				"settings",
				GetAppPath(),
				dir,
				viper.GetViper(),
				viper.New(),
				nil,
				nil,
			}
			settingConf.Init()
		} else {
			return nil
		}
	}
	return settingConf
}

// 用户输入 TestCase 的路径，本函数创建对应的路径
// GetUserDataPath()/testcase/<hashid>
// 并为其创建 <configtype>.yaml 的 viper
func CreateCacheConfiger(folder, configtype string) *Configer {
	if id, err := Hash(folder); err == nil {
		c := &Configer{
			configtype,
			"",
			GetSettingConfiger().GetString("usercachepath") + "/" + id,
			viper.New(),
			viper.New(),
			nil,
			nil,
		}
		c.Init()
		return c
	}
	return nil
}

func ReadYaml(yamldir string) interface{} {
	if _, err := os.Stat(yamldir); err != nil {
		log.Warn().Msg("Cannot read file") // 这不算 error
		return nil
	}
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName("config")
	v.AddConfigPath(yamldir)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 找不到 config file
			log.Error().Stack().Err(errors.Wrap(err, "Can not find config file")).Send()
		} else {
			// 可以找到 config file，但出了其他 error
			log.Error().Stack().Err(errors.Wrap(err, "Found config file, but other error accord")).Send()
		}
	} else {
		return v.AllSettings()
	}
	return nil
}
