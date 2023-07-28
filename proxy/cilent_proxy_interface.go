package proxy

type ClientProxy interface {
	Invoke(methodName string, args []interface{}) interface{}
}
