package factory

import (
	"github.com/fyupeng/rpc-go-netty/net/netty/future"
	"github.com/modern-go/concurrent"
	"log"
	"reflect"
)

type singleton struct {
	Data concurrent.Map
	// 这里可以添加其他需要的字段
}

var instance = &singleton{}

func GetInstance(structType reflect.Type) interface{} {
	var structValue interface{}
	switch structType.Elem().Name() {
	case "UnProcessResult":
		var hasFound bool
		structValue, hasFound = instance.Data.Load(structType)
		if hasFound {
			return structValue
		}
		structValue = future.NewUnprocessResult()
		instance.Data.Store(structType, structValue)
	default:
		log.Fatal("GetInstance failed: not found struct Type ", structType)
	}

	return structValue
}
