package imes

import (
	"fmt"
	"testing"
)

// TestXXX() :go test // 会自动调用并报告结果 PASS/FAIL
func TestAddCounter(t *testing.T) {
	ta := new(Api)
	ta.InitCounter()
	i := ta.GetCounter()
	if (i + 1) != ta.AddCounter(1) {
		t.Errorf("AddCounter: want %d, got %d", i+1, i)
	}
	if (i + +1 + 100) != ta.AddCounter(100) {
		t.Errorf("AddCounter: want %d, got %d", i+100, i)
	}
}

// BenchmarkXXX(): go test 会多次执行并计算一个平均执行时间
func BenchmarkAddCounter(b *testing.B) {
	ta := new(Api)
	for i := 0; i < b.N; i++ {
		ta.AddCounter(1)
	}
}

func BenchmarkAddCounterParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		ta := new(Api)
		for i := 0; i < b.N; i++ {
			ta.AddCounter(1)
		}
	})
}

// Example<Type>_<Func><opt>() : 自动测试 + 生成文档，注释中没有 Output: 不予执行
func ExampleApi_AddCounter() {
	ta := new(Api)
	ta.InitCounter()
	fmt.Println(ta.AddCounter(1))
	// output:
	// 1
}

func TestFileWalk(t *testing.T) {
	a := new(Api)
	fs, err := a.WalkMatch(GetAppPath(), "*.py")
	if err != nil {
		panic(err)
	}
	t.Log(fs)
}
