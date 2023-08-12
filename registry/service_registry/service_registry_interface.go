package service_registry

/*
*
服务注册器接口
*/
type ServiceRegistry interface {
	AddService(service interface{}, serviceName string)

	Listen() error

	Register(serviceName string) (err error)

	RegisterWithGroupName(serviceName string, groupName string) (err error)

	ClearRegistry() (err error)
}
