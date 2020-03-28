package config

import (
	"context"
)

type MockSD struct {
	Address string
	Port    int
	Err     error
}

func (m MockSD) GetService(ctx context.Context, name string) (string, int, error) {
	return m.Address, m.Port, m.Err
}
