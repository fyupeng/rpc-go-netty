package aop

import (
	"bou.ke/monkey"
	"fmt"
	"reflect"
	"regexp"
)

var aspectList = make([]AspectInterface, 0)

func RegisterPoint(pointType reflect.Type) {
	pkgPath := pointType.PkgPath()
	receiverName := pointType.Name()
	// pointType 为 指针类型,获取指针的基础类型(例如：*int -> int)
	if pointType.Kind() == reflect.Ptr {
		pkgPath = pointType.Elem().PkgPath()
		receiverName = pointType.Elem().Name()
	}
	// 遍历该结构具有的函数
	for i := 0; i < pointType.NumMethod(); i++ {
		method := pointType.Method(i)
		// 方法位置字符串 “报名.接受者.方法名”，用于匹配代理
		methodLocation := fmt.Sprintf("%s.%s.%s", pkgPath, receiverName, method.Name)
		fmt.Println(methodLocation)
		var guard *monkey.PatchGuard
		//func(*aop.Hello, string, string) string -> in []reflect.Value -> [*aop.Hello, string, string]
		var proxy = func(in []reflect.Value) []reflect.Value {
			guard.Unpatch()
			defer guard.Restore()
			receiver := in[0]
			point := NewJoinPoint(receiver, in[1:], method)
			fmt.Println("methodLocation", methodLocation)
			defer finallyProcessed(point, methodLocation)
			// 执行 前置处理
			if !beforeProcessed(point, methodLocation) {
				return point.Result
			}
			for i := 0; i < method.Func.Type().NumOut(); i++ {
				fmt.Printf("%s,", method.Func.Type().Out(i))
			}

			// 执行用户函数
			point.Result = receiver.MethodByName(method.Name).Call(in[1:])
			// 后置处理
			afterProcessed(point, methodLocation)
			return point.Result
		}
		// 动态创建代理函数
		proxyFn := reflect.MakeFunc(method.Func.Type(), proxy)
		// 利用 monkey 框架替换 被代理函数
		guard = monkey.PatchInstanceMethod(pointType, method.Name, proxyFn.Interface())
	}

}

func RegisterAspect(aspect AspectInterface) {
	aspectList = append(aspectList, aspect)
}

// 前置处理
func beforeProcessed(point *JoinPoint, methodLocation string) bool {
	for _, aspect := range aspectList {
		if !isAspectMatch(aspect.GetAspectExpress(), methodLocation) {
			continue
		}
		if !aspect.Before(point) {
			return false
		}
	}
	return true
}

// 后置处理
func afterProcessed(point *JoinPoint, methodLocation string) {
	for i := len(aspectList) - 1; i >= 0; i-- {
		aspect := aspectList[i]
		if !isAspectMatch(aspect.GetAspectExpress(), methodLocation) {
			continue
		}
		aspect.After(point)
	}
}

// 最终处理
func finallyProcessed(point *JoinPoint, methodLocation string) {
	for i := len(aspectList) - 1; i >= 0; i-- {
		aspect := aspectList[i]
		if !isAspectMatch(aspect.GetAspectExpress(), methodLocation) {
			continue
		}
		aspect.Finally(point)
	}
}

func isAspectMatch(aspectExpress, methodLocation string) bool {
	// aspectExpress 采用正则表达式
	pattern, err := regexp.Compile(aspectExpress)
	if err != nil {
		return false
	}
	return pattern.MatchString(methodLocation)
}
