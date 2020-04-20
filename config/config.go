// Package config is responsible for loading user-manager application config.
// Basic configuration like consul credentials and address, http port to listen for requests,
// postgres schema name, credentials, and client timeout are read from environment variables.
package config

import (
	"context"
	"fmt"
	"time"

	"github.com/lvl484/user-manager/logger"
	"github.com/lvl484/user-manager/server/mail"
	"github.com/lvl484/user-manager/storage"

	consul "github.com/hashicorp/consul/api"
	"github.com/kelseyhightower/envconfig"
)

// Config model includes all necessary information, which will be read from environment variables
type Config struct {
	PostgresUser string `envconfig:"POSTGRES_USER" required:"true"`
	PostgresPass string `envconfig:"POSTGRES_PASSWORD" required:"true"`
	PostgresDB   string `envconfig:"POSTGRES_DB" required:"true"`

	ConsulAddress string `envconfig:"CONSUL_ADDRESS" required:"true"`
	ConsulToken   string `envconfig:"CONSUL_TOKEN"`

	HTTPIP       string        `envconfig:"HTTP_IP" default:"0.0.0.0"`
	HTTPPort     int           `envconfig:"HTTP_PORT" default:"8000"`
	ReadTimeout  time.Duration `envconfig:"READ_TIMEOUT" default:"60s"`
	WriteTimeout time.Duration `envconfig:"WRITE_TIMEOUT" default:"60s"`

	LoggerPassSecret string `envconfig:"LOGGER_PASS_SECRET"`
	LoggerPassSHA2   string `envconfig:"LOGGER_PASS_SHA2"`
	LoggerOutput     string `envconfig:"LOGGER_OUTPUT" default:"Stdout"`
	LoggerLevel      string `envconfig:"LOGGER_LEVEL" default:"info"`
	LoggerType       string `envconfig:"LOGGER_TYPE" default:"async"`

	EmailAddress  string `envconfig:"EMAIL_ADDRESS" default:"user.namager@gmail.com"`
	EmailPassword string `envconfig:"EMAIL_PASSWORD" default:"lvl484Golang"`
	EmailHost     string `envconfig:"EMAIL_HOST" default:"smtp.gmail.com"`
	EmailPort     int    `envconfig:"EMAIL_PORT" default:"587"`

	sd ServiceDiscovery
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
	consulClient, err := consul.NewClient(consulConfig)
	if err != nil {
		return nil, fmt.Errorf("consul client error %w", err)
	}

	config.sd = consulSD{consul: consulClient}

	return &config, nil
}

// LoggerConfig get configurations for glaylog
func (c *Config) LoggerConfig(ctx context.Context) (*logger.LogConfig, error) {
	const serviceName = "graylog"

	host, port, err := c.sd.GetService(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	return &logger.LogConfig{
		Host:       host,
		Port:       port,
		PassSecret: c.LoggerPassSecret,
		PassSHA2:   c.LoggerPassSHA2,
		Output:     c.LoggerOutput,
		Level:      c.LoggerLevel,
		Type:       c.LoggerType,
	}, nil
}

// DBConfig get configuration for Postgres Database
func (c *Config) DBConfig(ctx context.Context) (*storage.DBConfig, error) {
	const serviceName = "db"

	host, port, err := c.sd.GetService(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	return &storage.DBConfig{
		Host:     host,
		Port:     port,
		User:     c.PostgresUser,
		Password: c.PostgresPass,
		DBName:   c.PostgresDB,
	}, nil
}

func (c *Config) ServerAddress() string {
	return fmt.Sprintf("%s:%d", c.HTTPIP, c.HTTPPort)
}

// EmailConfig get configuration for Email
func (c *Config) EmailConfig() (*mail.EmailInfo, error) {
	return &mail.EmailInfo{
		Sender:   c.EmailAddress,
		Password: c.EmailPassword,
		Host:     c.EmailHost,
		Port:     c.EmailPort,
	}, nil
}
