package load_balancer

import (
	"errors"
	"math/rand"
	"reflect"
	"time"
)

type LoadBalancer struct {
}

func (loadBalancer *LoadBalancer) SelectService(services interface{}) (interface{}, error) {
	servicesVal := reflect.ValueOf(services)

	if servicesVal.Kind() != reflect.Slice {
		return nil, errors.New("services is not a slice")
	}

	if servicesVal.Len() == 0 {
		return nil, errors.New("services is empty")
	}

	//if len(services) == 0 {
	//	var empty interface{}
	//	return empty, errors.New("loadBalancer can't discover service!")
	//}
	rand.Seed(time.Now().Unix())
	//return services[rand.Int()], nil
	return servicesVal.Index(rand.Intn(servicesVal.Len())).Interface(), nil
}

func (loadBalancer *LoadBalancer) SelectNode(nodes []string) (string, error) {
	if len(nodes) == 0 {
		return "", errors.New("loadBalancer can't discover service!")
	}
	rand.Seed(time.Now().Unix())
	return nodes[rand.Int()], nil
}
