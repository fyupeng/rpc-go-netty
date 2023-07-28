package aop

import (
	"reflect"
)

// 连接点
type JoinPoint struct {
	Receiver interface{}
	Method   reflect.Method
	Params   []reflect.Value
	Result   []reflect.Value
}

/*
*
receiver:
*/
func NewJoinPoint(receiver interface{}, params []reflect.Value, method reflect.Method) *JoinPoint {
	point := &JoinPoint{
		Receiver: receiver,
		Params:   params,
		Method:   method,
	}
	//
	fn := method.Func
	fnType := fn.Type()
	nout := fnType.NumOut()
	point.Result = make([]reflect.Value, nout)

	for i := 0; i < nout; i++ {
		// 对函数的返回值列表元素值进行初始化，默认为返回空值
		point.Result[i] = reflect.Zero(fnType.Out(i))
	}
	return point
}
