package testset

// 测试项日志
type TestItemLog struct {
	TestItemId int    `json:"testItemId"`
	Message    string `json:"message"`
	TimeStamp  int64  `json:"timestamp"`
}
