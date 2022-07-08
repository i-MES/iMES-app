package target

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/i-mes/imes-app/backend/utils"
	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

/*
TestGroup
数据I/O:
	1. Parse: 从 python、go 等文件中解析(2种策略：所有文件1个Group、每个文件1个Group)
	2. Save/Read: 写入本地 testgroup.json 文件，或从该文件中读取。—— 用户 UI 上的修改实时写入 testgroup.json。
	3. Down/UpLoad: 从网络上加载/上传 —— testgroup.json 定时 upload 到云端。
数据冲突问题:
	首先应该尽量避免出现冲突，但多种途径的数据IO一旦出现冲突，还是不可避免的要解决冲突。场景及策略：
	1. 源代码修改，parse 到新的 tc 和 tg：
	2.
*/
type TestGroup struct {
	Id          string      `json:"id"`
	Title       string      `json:"title"`
	Desc        string      `json:"desc"`
	TestClasses []TestClass `json:"testclasses"`
}

// strategy(策略)
// 1. 所有 py 文件生成 1 个 TestGroup
// 2. 每个 py 文件生成 1 个 TestGroup
func ParsePythons(ctx *context.Context, filepaths []string, strategy int) []TestGroup {
	tgs := make([]TestGroup, 0)
	if strategy == 1 {
		if _uuid, err := uuid.NewUUID(); err == nil {
			tg := TestGroup{_uuid.String(), "", "", make([]TestClass, 0)}
			for _, fp := range filepaths {
				tg.TestClasses = append(tg.TestClasses, ParsePythonOneStep(ctx, fp)...)
			}
			tgs = append(tgs, tg)
		} else {
			fmt.Println("error uuid get: ", err)
		}
	} else if strategy == 2 {
		for _, fp := range filepaths {
			if _uuid, err := uuid.NewUUID(); err == nil {
				tgs = append(tgs, TestGroup{_uuid.String(), path.Base(fp), fp, ParsePythonOneStep(ctx, fp)})
			} else {
				fmt.Println("error uuid get: ", err)
			}
		}
	}
	wails.LogDebug(*ctx, "ParsePythons tgs len:"+string(len(tgs)))
	return tgs
}

// 从源代码文件文件中加载 TG 数据，并开启 SyncMonitor
func LoadTestGroupFromSrc(ctx *context.Context, selectPath bool) ([]TestGroup, []string) {
	folderpath := ""
	if selectPath {
		// 用户选择文件夹
		folderpath = utils.OpenFolder(ctx, "Open TestCase Folder")
	} else {
		// 使用默认文件夹
		folderpath = utils.GetAppPath() + "/testcase/python/"
	}
	filepathes, err := utils.GetAllFile(folderpath, "test*.py", true)
	if err != nil {
		wails.LogDebug(*ctx, "Can not load python file")
		return nil, nil
	}
	tgs := ParsePythons(ctx, filepathes, 2)
	return tgs, filepathes
}

var tgMonitor *time.Ticker

// testgroupMonitorHandler.Stop() 能够停止 Monitor

// 从 Config 文件中加载 TG 数据，并开启 SyncMonitor
func LoadTestGroupFromConfig(ctx *context.Context) ([]TestGroup, []string) {
	ctgs := make([]TestGroup, 0)

	if data, err := utils.InputConfigData("testgroup"); err == nil {
		// 找到 json 文件，加载到 ctgs
		_tg := []byte(json.Get(data).ToString())
		if json.Unmarshal(_tg, &ctgs) != nil {
			fmt.Println("can not Unmarshal json data")
			return nil, nil
		}
		srcs := make([]string, 0)
		for _, tg := range ctgs {
			for _, tc := range tg.TestClasses {
				srcs = append(srcs, tc.FileName)
			}
		}
		return ctgs, srcs
	} else if os.IsNotExist(err) {
		wails.LogDebug(*ctx, "json 文件不存在")
		return nil, nil
	}
	return nil, nil
}

/*
加载 TestGroup 到内存（Cache）

加载源 & 加载方式：
	loadConfig: 从 config/xxx.json 中加载数据
	loadFile:   从 源代码文件中解析数据，并完整重写 Config 文件
		selectPath: 源代码文件的目录由用户选择 or App 默认目录

同步策略：
	- 2 种加载方式都会启动 SyncMonitor
	- Cache 实时向 Config 同步：用户在 UI 上的修改（拖拽等）会实时保存到 Config 文件，但只修改 TG，不修改 TC。
	- 源文件监控到修改后向 Cache 同步，并同时写入 Config。

所以：
	用户仅需解析一次源文件，后续源文件修改会自动同步 Merge 到 Cache & Config。
	用户再次解析源文件，也要确认是否已存在 Config，存在则进行 Merge。
*/
func LoadTestGroup(ctx *context.Context, loadFlag string, selectPath bool) []TestGroup {
	switch loadFlag {
	case "config":
		ctgs, fs := LoadTestGroupFromConfig(ctx)
		StartTestGroupSyncMonitor(ctx, fs, true)
		return ctgs
	case "src":
		stgs, srcs := LoadTestGroupFromSrc(ctx, selectPath)
		StartTestGroupSyncMonitor(ctx, srcs, true)
		SaveTestGroup(ctx, stgs)
		return stgs
	}
	return nil
}

func StartTestGroupSyncMonitor(ctx *context.Context, srcs []string, autoMerge bool) {
	StopTestGroupSyncMonitor() // 只使用1个TG Monitor，先删除原有的。
	wails.EventsEmit(*ctx, "testgroupmonitor", "clear")
	tgMonitor = utils.StartSyncMonitor("testgroup", srcs,
		func(newsrcfile string) {
			fmt.Println("---------------")
			if autoMerge {
				// 从源文件读取 TG ，并自动 Merge
				ctgs, _ := LoadTestGroupFromConfig(ctx)
				stgs := ParsePythons(ctx, []string{newsrcfile}, 1)
				for _, stc := range stgs[0].TestClasses {
					// Merge 策略：
					//   src 中新增的 TC：保留到新的 TG 中，添加到 ctgs；
					//   src 中删除的 TC：
					//   src 中修改的 TC：
					//   src 中新增的 TI：保留到新的 TG 中，添加到 ctgs；
					//   src 中删除的 TI：
					//   src 中修改的 TI：
					for i, ctg := range ctgs {
						for j, ctc := range ctg.TestClasses {
							if ctc.ClassName == stc.ClassName {
								ctgs[i].TestClasses[j].TestItems = stc.TestItems
								break
							}
						}
					}
					SaveTestGroup(ctx, ctgs)
					wails.LogInfo(*ctx, "Merge and save newer testgroup")
				}
			}
			wails.EventsEmit(*ctx, "testgroupmonitor", "srcnewer")
		},
		func() {
			fmt.Println("===============")
			wails.EventsEmit(*ctx, "testgroupmonitor", "srcolder")
		},
	)
}
func StopTestGroupSyncMonitor() {
	if tgMonitor != nil {
		tgMonitor.Stop()
	}
}

func SaveTestGroup(ctx *context.Context, data []TestGroup) {
	_data := make(map[string]interface{})
	_data["testgroup"] = data
	utils.OutputConfigData(_data)
}

// 保存 TestGroup 信息
func (tg *TestGroup) Save(ctx context.Context) {

}

// 在本地保存信息、网络保存信息、源码文件中寻找差异，并进行 merge
func (tg *TestGroup) Merge(ctx context.Context) {

}

// 执行 TestGroup 内 TestClass 的测试，Group 内串行，Group 间并行
func (tg *TestGroup) Run(ctx *context.Context, teid string) {
	for _, tc := range tg.TestClasses {
		tc.Run(ctx, func(ename string, tiid string, msg string) {
			wails.EventsEmit(*ctx, ename,
				TestItemStatus{teid, tg.Id, tiid, time.Now().Unix(), msg})
		})
	}
	wails.EventsEmit(*ctx, "testgroupfinished", tg.Id)
}

func (tg *TestGroup) Pause(ctx *context.Context) {
}

func (tg *TestGroup) Stop(ctx *context.Context) {
}
