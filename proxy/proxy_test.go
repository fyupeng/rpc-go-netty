package proxy

import (
	"context"
	"fmt"
	"time"
)

// 装饰器函数类型
type DecoratorFunc func(context.Context, func()) error

// 装饰器函数
func proxyDecorator(ctx context.Context, fn func()) error {
	// 在调用目标函数之前执行一些逻辑
	fmt.Println("Before executing the function")

	// 执行目标函数
	fn()

	// 在调用目标函数之后执行一些逻辑
	fmt.Println("After executing the function")

	return nil
}

// 动态代理函数
func Proxy(fn interface{}, decorator DecoratorFunc) interface{} {
	return func(ctx context.Context) error {
		if fn, ok := fn.(func()); ok {
			return decorator(ctx, fn)
		}
		return fmt.Errorf("invalid function type")
	}
}

// 要被代理的函数
func TargetFunction() {
	fmt.Println("Executing target function...")
	time.Sleep(1 * time.Second)
	fmt.Println("Target function execution completed.")
}
