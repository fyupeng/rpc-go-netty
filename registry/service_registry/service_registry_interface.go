package service_registry

/*
*
服务注册器接口
*/
type ServiceRegistry interface {
	Listen() error

	Register(serviceName string) (err error)

	RegisterWithGroupName(serviceName string, groupName string) (err error)

	ClearRegistry() (err error)
}
