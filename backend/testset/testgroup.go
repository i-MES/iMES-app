package testset

// 测试组
type TestGroup struct {
	Id        int        `json:"id"`
	Title     string     `json:"title"`
	Desc      string     `json:"desc"`
	TestItems []TestItem `json:"testItems"`
}
