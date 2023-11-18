package factory

import "sync"

type Singleton struct {
	Data interface{}
	// 这里可以添加其他需要的字段
}

var instance *Singleton
var once sync.Once

func GetInstance(message interface{}) *Singleton {
	once.Do(func() {
		instance = &Singleton{
			Data: message, // 设置实例的初始数据
		}
	})
	return instance
}
