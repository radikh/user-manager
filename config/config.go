// Package config is responsible for loading user-manager application config.
// Basic configuration like consul credentials and address, http port to listen for requests,
// postgres schema name, credentials, and client timeout are read from environment variables.
package config

import (
	"context"
	"errors"
	"fmt"
	"time"

	consul "github.com/hashicorp/consul/api"
	"github.com/kelseyhightower/envconfig"
)

// Config model includes all necessary information, which will be read from environment variables
type Config struct {
	PostgresUser string `envconfig:"POSTGRES_USER" required:"true"`
	PostgresPass string `envconfig:"POSTGRES_PASSWORD" required:"true"`
	PostgresDB   string `envconfig:"POSTGRES_DB" required:"true"`

	ConsulAddress string `envconfig:"CONSUL_ADDRESS" required:"true"`
	ConsulToken   string `envconfig:"CONSUL_TOKEN" required:"true"`

	BindIP   string `envconfig:"BIND_IP" default:"0.0.0.0"`
	BindPort int    `envconfig:"BIND_PORT" default:"8000"`

	Timeout time.Duration `envconfig:"TIMEOUT" default:"60s"`

	consulClient *consul.Client
}

// NewConfig() create new configuration for application
func NewConfig() (*Config, error) {
	var config Config
	// Read all environment variables and fill config structure with them
	err := envconfig.Process("", &config)
	if err != nil {
		return nil, fmt.Errorf("envconfig error %w", err)
	}

	// initialization configuration for consul client
	consulConfig := &consul.Config{
		Address: config.ConsulAddress,
		Token:   config.ConsulToken,
	}

	// Create new consul client using prepared configuration
	config.consulClient, err = consul.NewClient(consulConfig)
	if err != nil {
		return nil, fmt.Errorf("consul client error %w", err)
	}

	return &config, nil
}

// Get configurations for glaylog
func (c *Config) GraylogConfig(ctx context.Context) (string, error) {
	const serviceName = "graylog"

	opts := new(consul.QueryOptions).WithContext(ctx)

	services, _, err := c.consulClient.Catalog().Service(serviceName, "", opts)
	if err != nil {
		return "", fmt.Errorf("resolve graylog service error %w", err)
	}

	if len(services) == 0 {
		return "", errors.New("graylog service not found")
	}

	host := services[0].Address
	port := services[0].ServicePort

	return fmt.Sprintf("%s:%d", host, port), nil
}

// Get configuration for Postgres Database
func (c *Config) PostgresConfig(ctx context.Context) (string, error) {
	const serviceName = "db"

	opts := new(consul.QueryOptions).WithContext(ctx)

	services, _, err := c.consulClient.Catalog().Service(serviceName, "", opts)
	if err != nil {
		return "", fmt.Errorf("resolve db service error %w", err)
	}

	if len(services) == 0 {
		return "", errors.New("db service not found")
	}

	host := services[0].Address
	port := services[0].ServicePort

	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, c.PostgresUser, c.PostgresPass, c.PostgresDB), nil
}

func (c *Config) ServerAddress() string {
	return fmt.Sprintf("%s:%d", c.BindIP, c.BindPort)
}
