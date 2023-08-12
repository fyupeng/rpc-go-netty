package handler

/*
*
反射处理器接口
*/
type InvocationHandler interface {
	Handle(service interface{}, methodName string, args []interface{}) interface{}
}
