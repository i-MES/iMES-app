package utils

import (
	"context"

	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

/*
wails log level:

- Trace
- Debug
- Info
- Warning
- Error
- Fatal

设置 Log Level：
Go Signature: LogSetLogLevel(ctx context.Context, level logger.LogLevel)
JS Signature: LogSetLogLevel(level: number)
	level: number 取值：
		1	Trace
		2	Debug
		3	Info
		4	Warning
		5	Error

*/

var ctx context.Context

func InitLog(cont context.Context) {
	ctx = cont
}

func LogTrace(message string) {
	wails.LogTrace(ctx, message)
}
func LogDebug(message string) {
	wails.LogDebug(ctx, message)
}

func LogInfo(message string) {
	wails.LogInfo(ctx, message)
}

func LogWarning(message string) {
	wails.LogWarning(ctx, message)
}
func LogError(message string) {
	wails.LogError(ctx, message)
}
