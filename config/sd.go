package config

import (
	"context"
	"fmt"

	consul "github.com/hashicorp/consul/api"
)

type ServiceDiscovery interface {
	GetService(ctx context.Context, name string) (string, int, error)
}

type consulSD struct {
	consul *consul.Client
}

func (s consulSD) GetService(ctx context.Context, name string) (string, int, error) {
	opts := new(consul.QueryOptions).WithContext(ctx)

	services, _, err := s.consul.Catalog().Service(name, "", opts)
	if err != nil {
		return "", 0, fmt.Errorf("resolve %s service error %w", name, err)
	}

	if len(services) == 0 {
		return "", 0, fmt.Errorf("%s service not found", name)
	}

	return services[0].Address, services[0].ServicePort, nil
}
