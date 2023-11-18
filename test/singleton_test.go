package test

import (
	"fmt"
	"reflect"
	"rpc-go-netty/net/netty/future"
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
