package testset

import "fmt"

// 测试项
type TestItem struct {
	Title    string `json:"title"`
	Desc     string `json:"desc"`
	Funcname string `json:"funcname"`
	Sequence int    `json:"sequence"`
}

func (ti *TestItem) Run() {
	fmt.Println("Run TI: ", ti.Title, ti.Desc, ti.Funcname, ti.Sequence)
}
