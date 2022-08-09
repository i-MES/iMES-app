package imes

import (
	"fmt"
	"testing"
)

// TestXXX() :go test // 会自动调用并报告结果 PASS/FAIL
func TestUUID(t *testing.T) {
	a := new(Api)
	uuid := a.UUID()
	if len(uuid) != 36 {
		t.Errorf("UUID() return bits: want %d, got %d", 36, len(uuid))
	}
}

// BenchmarkXXX(): go test 会多次执行并计算一个平均执行时间
func BenchmarkUUID(b *testing.B) {
	ta := new(Api)
	for i := 0; i < b.N; i++ {
		ta.UUID()
	}
}

func BenchmarkAddUUIDParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		ta := new(Api)
		for i := 0; i < b.N; i++ {
			ta.UUID()
		}
	})
}

// Example<Type>_<Func><opt>() : 自动测试 + 生成文档，注释中没有 Output: 不予执行
func ExampleApi_UUID() {
	ta := new(Api)
	ta.UUID()
	fmt.Println(len(ta.UUID()))
	// output:
	// 36
}
