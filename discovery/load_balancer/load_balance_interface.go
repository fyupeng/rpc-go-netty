package load_balancer

type LoadBalanceInterface interface {
	SelectService(services interface{}) (interface{}, error)

	SelectNode(nodes []string) (string, error)
}
