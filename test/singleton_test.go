package test

import (
	"fmt"
	"github.com/fyupeng/rpc-go-netty/net/netty/future"
	"reflect"
	"testing"
)

func TestSingleton(t *testing.T) {
	typ := reflect.TypeOf((*future.UnProcessResult)(nil))
	fmt.Println(typ)
	fmt.Println(typ.Elem())
	fmt.Println(typ.Elem().Name())
	fmt.Println(typ.String())
	fmt.Println(typ.Elem().String())
}
