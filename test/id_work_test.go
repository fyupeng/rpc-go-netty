package test

import (
	"fmt"
	"github.com/fyupeng/rpc-go-netty/utils/idworker"
	"strconv"
	"testing"
)

func TestIdWorker(t *testing.T) {
	worker, err := idworker.NewNowWorker(0)
	if err != nil {
		fmt.Println(err)
	}
	for i := 0; i < 5; i++ {
		id := worker.NextId()
		fmt.Println(id)
		fmt.Println(strconv.FormatInt(id, 10))
	}

	// 结果：
	//180089216334561280
	//180089216334561281
	//180089216334561282
	//180089216334561283
	//180089216334561284

	//BenchmarkID-16           4902658（执行次数）              244.5 ns/op（平均每次执行所需时间）
	//BenchmarkID-4            3137922               		   366.1 ns/op
}

func BenchmarkID(b *testing.B) {
	worker, _ := idworker.NewNowWorker(0)
	for i := 0; i < b.N; i++ {
		worker.NextId()
	}
}
