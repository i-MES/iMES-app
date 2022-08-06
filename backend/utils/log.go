package utils

import (
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

/*
wails log 共 6 级 level:

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

zerolog 共 7 级 level
-	panic (zerolog.PanicLevel, 5)
-	fatal (zerolog.FatalLevel, 4)
-	error (zerolog.ErrorLevel, 3)
-	warn (zerolog.WarnLevel, 2)
-	info (zerolog.InfoLevel, 1)
-	debug (zerolog.DebugLevel, 0)
-	trace (zerolog.TraceLevel, -1)
*/

type AppLog struct {
	Level     string `json:"level"`
	TimeStamp int64  `json:"timestamp"`
	Source    string `json:"source"`
	Message   string `json:"message"`
}

func InitLog(level string) {
	zerolog.LevelFieldName = "l"
	zerolog.TimestampFieldName = "t"
	zerolog.MessageFieldName = "m"
	zerolog.ErrorFieldName = "e"
	zerolog.CallerFieldName = "c"
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix                 // 使用 Unix 时间戳
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack             // 开启栈跟踪
	zerolog.CallerMarshalFunc = func(file string, line int) string { // Caller 记录文件名不含路径
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}
	if level == "dev" {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// UTC（Unix标准时间）和 CST（中国时间）相差8个小时；
	// time.Now()输出默认 CST 时区时间；
	// time.Parse() 默认输出 UTC 时区时间。
	// Local() 得到当地时区
	time.LoadLocation(GetSettingConfiger().GetString("log.timezone"))
	log.Info().Msgf("Time: %v", time.Now().UTC().Local().Format(time.RFC3339))

	// 配置 Stdout 的输出
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC850}
	consoleWriter.FormatTimestamp = func(i interface{}) string {
		if v, ok := i.(int64); ok {
			return time.Unix(v, 0).Format(time.RFC3339)
		} else {
			return time.Now().UTC().Local().Format(time.RFC3339)
		}
	}

	fileWriter, _ := os.OpenFile("./.logs/log.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	// 配置文件及 rotate
	rotateFileWriter, err := NewRotate(RotateOptions{
		Directory:       GetSettingConfiger().GetString("log.logsfilepath"),
		MaximumFileSize: 1024 * 1024 * 1024,
		MaximumLifetime: time.Minute * time.Duration(GetSettingConfiger().GetInt("log.logsfilerotateminutes")),
		FileNameFunc:    nil,
	})

	if err == nil {
		// defer rotateFileWriter.Close() 不能关闭啊，那还记录个啥？

		multi := zerolog.MultiLevelWriter(consoleWriter, rotateFileWriter, fileWriter)
		log.Logger = zerolog.New(multi).With().Timestamp().
			Caller(). // 默认添加 Caller
			Stack().  // Stack 只会在 Err() 时有效
			Logger()
	}
}
