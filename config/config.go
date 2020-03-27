// Package config is responsible for loading user-manager application config.

package config

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	consul "github.com/hashicorp/consul/api"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"

	"github.com/lvl484/user-manager/logger"
)

const (
	loggerHost       = "loggerUM.Host"
	loggerPort       = "loggerUM.Port"
	loggerPassSecret = "loggerUM.PassSecret"
	loggerPassSHA2   = "loggerUM.PassSHA2"
	loggerOutput     = "loggerUM.Output"
	LoggerLevel      = "loggerUM.Level"
	LoggerType       = "loggerUM.Type"
)

// Config model includes all necessary information, which will be read from environment variables
type Config struct {
	PostgresUser string `envconfig:"POSTGRES_USER" required:"true"`
	PostgresPass string `envconfig:"POSTGRES_PASSWORD" required:"true"`
	PostgresDB   string `envconfig:"POSTGRES_DB" required:"true"`

	ConsulAddress string `envconfig:"CONSUL_ADDRESS" required:"true"`
	ConsulToken   string `envconfig:"CONSUL_TOKEN" required:"true"`

	BindIP       string        `envconfig:"BIND_IP" default:"0.0.0.0"`
	BindPort     int           `envconfig:"BIND_PORT" default:"8000"`
	ReadTimeout  time.Duration `envconfig:"READ_TIMEOUT" default:"60s"`
	WriteTimeout time.Duration `envconfig:"WRITE_TIMEOUT" default:"60s"`

	consulClient *consul.Client
	v            *viper.Viper
}

func NewViperConfig(configName, configPath string) (*Config, error) {
	v := viper.New()
	v.AddConfigPath(configPath)
	v.SetConfigName(configName)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	return &Config{v: v}, nil
}

// NewPostgresConfig returns pointer to PointerConfig with data read from viper.config.json
func (conf *Config) NewLoggerConfig() *logger.LogConfig {
	return &logger.LogConfig{
		Host:       conf.v.GetString(loggerHost),
		Port:       conf.v.GetString(loggerPort),
		PassSecret: conf.v.GetString(loggerPassSecret),
		PassSHA2:   conf.v.GetString(loggerPassSHA2),
		Output:     conf.v.GetString(loggerOutput),
		Level:      conf.v.GetString(LoggerLevel),
		Type:       conf.v.GetString(LoggerType),
	}
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

// LoggerConfig get configurations for glaylog
func (c *Config) LoggerConfig(ctx context.Context) (string, error) {
	const serviceName = "graylog"
	c.ServerAddress()
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

// DBConfig get configuration for Postgres Database
func (c *Config) DBConfig(ctx context.Context) (string, error) {
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
