package testset

import (
	"context"
	"time"

	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

// 测试项日志
type TestItemLog struct {
	TestItemId int    `json:"testItemId"`
	Flag       bool   `json:"flag"`
	Message    string `json:"message"`
	TimeStamp  int64  `json:"timestamp"`
}

func EmitTestItemLog(ctx context.Context, flag bool, message string) {
	wails.EventsEmit(ctx, "testitemlog",
		TestItemLog{1, flag, message, time.Now().Unix()})
}
