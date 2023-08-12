package load_balancer

import (
	"errors"
	"reflect"
)

func NewRoundRobinLoadBalancer() LoadBalancer {
	return &roundRobinLoadBalancer{
		Index: 0,
	}
}

type roundRobinLoadBalancer struct {
	Index int
}

func (loadBalancer *roundRobinLoadBalancer) SelectService(services interface{}) (interface{}, error) {
	servicesVal := reflect.ValueOf(services)

	if servicesVal.Kind() != reflect.Slice {
		return nil, errors.New("services is not a slice")
	}

	if servicesVal.Len() == 0 {
		return nil, errors.New("services is empty")
	}

	//if len(services) == 0 {
	//	var empty interface{}
	//	return empty, errors.New("loadBalancer can't discover cn.fyupeng.service!")
	//}
	//return services[rand.Int()], nil
	loadBalancer.Index += 1
	return servicesVal.Index(loadBalancer.Index % servicesVal.Len()).Interface(), nil
}

func (loadBalancer *roundRobinLoadBalancer) SelectNode(nodes []string) (string, error) {
	if len(nodes) == 0 {
		return "", errors.New("loadBalancer can't discover cn.fyupeng.service!")
	}
	loadBalancer.Index += 1
	return nodes[loadBalancer.Index%len(nodes)], nil
}
