package imes

import (
	"context"

	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

// 前后端传递的通用信息
type IMESInfo struct {
	Source    string `json:"source"`   // 发送者，如: python module
	Function  string `json:"function"` // 发送函数
	Type      string `json:"type"`     // 类型: error、warnning、info、debug
	Message   string `json:"message"`
	TimeStamp int64  `json:"timestamp"`
}

func (i *IMESInfo) Emit(ctx context.Context) {
	wails.EventsEmit(ctx, "imesinfo", i)
}
