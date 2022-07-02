package target

import (
	"context"
	"fmt"
	"path"
	"time"

	"github.com/google/uuid"
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
func ParsePythons(filepaths []string, strategy int) []TestGroup {
	tgs := make([]TestGroup, 0)
	if strategy == 1 {
		if _uuid, err := uuid.NewUUID(); err == nil {
			tg := TestGroup{_uuid.String(), "", "", make([]TestClass, 0)}
			for _, fp := range filepaths {
				tg.TestClasses = append(tg.TestClasses, ParsePython(fp)...)
			}
			tgs = append(tgs, tg)
		} else {
			fmt.Println("error uuid get: ", err)
		}
	} else if strategy == 2 {
		for _, fp := range filepaths {
			if _uuid, err := uuid.NewUUID(); err == nil {
				tgs = append(tgs, TestGroup{_uuid.String(), path.Base(fp), fp, ParsePython(fp)})
			} else {
				fmt.Println("error uuid get: ", err)
			}
		}
	}
	fmt.Println("tgs len:", len(tgs))
	return tgs
}

// 保存 TestGroup 信息
func (tg *TestGroup) Save(ctx context.Context) {

}

// 在本地保存信息、网络保存信息、源码文件中寻找差异，并进行 merge
func (tg *TestGroup) Merge(ctx context.Context) {

}

// 执行 TestGroup 内 TestClass 的测试，Group 内串行，Group 间并行
func (tg *TestGroup) Run(ctx context.Context, teid string) {
	for _, tc := range tg.TestClasses {
		tc.Run(ctx, func(ename string, tiid string, msg string) {
			wails.EventsEmit(ctx, ename,
				TestItemStatus{teid, tg.Id, tiid, time.Now().Unix(), msg})
		})
	}
	wails.EventsEmit(ctx, "testgroupfinished", tg.Id)
}

func (tg *TestGroup) Pause(ctx context.Context) {
}

func (tg *TestGroup) Stop(ctx context.Context) {
}
