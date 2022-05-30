package testset

// 测试项
type TestItem struct {
	Title    string `json:"title"`
	Desc     string `json:"desc"`
	Funcname string `json:"funcname"`
	Sequence int    `json:"sequence"`
}
