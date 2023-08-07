package load_balancer

type LoadBalancer interface {
	SelectService(services interface{}) (interface{}, error)

	SelectNode(nodes []string) (string, error)
}
