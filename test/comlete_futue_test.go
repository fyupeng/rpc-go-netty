package test

import (
	"fmt"
	"github.com/fyupeng/rpc-go-netty/net/netty/future"
	"testing"
	"time"
)

func TestCompleteFuture(t *testing.T) {
	completeFuture := future.NewCompleteFuture(make(chan interface{}), time.Second*10)

	unprocessResult := future.NewUnprocessResult()
	unprocessResult.Put("123", completeFuture)

	go test(unprocessResult)

	fmt.Println(time.Now().String())
	future1, err1 := completeFuture.GetFuture()
	fmt.Println(time.Now().String())
	fmt.Println(err1)
	fmt.Println(future1)

}

func test(unprocessResult *future.UnProcessResult) {
	time.Sleep(time.Second * 5)
	unprocessResult.Complete("123", "hello complete result")
}
