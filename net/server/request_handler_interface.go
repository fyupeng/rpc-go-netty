package server

/*
*
反射处理器接口
*/
type InvocationHandler interface {
	Handle(methodName string, args []interface{}) interface{}
}
