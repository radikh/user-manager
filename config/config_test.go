package config

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfigRequired(t *testing.T) {
	os.Setenv("POSTGRES_USER", "postgres")
	os.Setenv("POSTGRES_PASSWORD", "1q2w3e4r")
	os.Setenv("POSTGRES_DB", "um_db")
	os.Setenv("CONSUL_ADDRESS", "consul:8500")
	os.Setenv("CONSUL_TOKEN", "token")

	cfg, err := NewConfig()
	require.NoError(t, err)

	testCases := []struct {
		name     string
		got      string
		expected string
	}{
		{
			name:     "POSTGRES_USER",
			got:      cfg.PostgresUser,
			expected: "postgres",
		}, {
			name:     "POSTGRES_PASSWORD",
			got:      cfg.PostgresPass,
			expected: "1q2w3e4r",
		}, {
			name:     "POSTGRES_DB",
			got:      cfg.PostgresDB,
			expected: "um_db",
		}, {
			name:     "CONSUL_ADDRESS",
			got:      cfg.ConsulAddress,
			expected: "consul:8500",
		}, {
			name:     "CONSUL_TOKEN",
			got:      cfg.ConsulToken,
			expected: "token",
		},
	}

	for _, tt := range testCases {
		assert.Equal(t, tt.got, tt.expected, tt.name)
	}
}

func TestNewConfigDefault(t *testing.T) {
	os.Setenv("POSTGRES_USER", "unused")
	os.Setenv("POSTGRES_PASSWORD", "unused")
	os.Setenv("POSTGRES_DB", "unused")
	os.Setenv("CONSUL_ADDRESS", "unused:8500")
	os.Setenv("CONSUL_TOKEN", "unused")

	cfg, err := NewConfig()
	require.NoError(t, err)

	testCases := []struct {
		name     string
		got      interface{}
		expected interface{}
	}{
		{
			name:     "HTTP_IP",
			got:      cfg.HTTPIP,
			expected: "0.0.0.0",
		}, {
			name:     "HTTP_PORT",
			got:      cfg.HTTPPort,
			expected: 8000,
		}, {
			name:     "READ_TIMEOUT",
			got:      cfg.ReadTimeout.Seconds(),
			expected: 60,
		}, {
			name:     "WRITE_TIMEOUT",
			got:      cfg.WriteTimeout.Seconds(),
			expected: 60,
		},
	}

	for _, tt := range testCases {
		assert.EqualValues(t, tt.got, tt.expected, tt.name)
	}
}

func TestNewConfig(t *testing.T) {
	os.Unsetenv("POSTGRES_USER")

	cfg, err := NewConfig()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "POSTGRES_USER")
	assert.Nil(t, cfg)
}

func TestConfigServerAddress(t *testing.T) {
	c := Config{
		HTTPIP:   "sun",
		HTTPPort: 9999,
	}
	assert.Equal(t, "sun:9999", c.ServerAddress())
}

func TestConfigDBConfig(t *testing.T) {
	sd := &MockSD{
		Address: "moon",
		Port:    9999,
	}
	c := Config{
		PostgresUser: "postgres",
		PostgresPass: "1q2w3e4r",
		PostgresDB:   "um_db",
		sd:           sd,
	}

	got, err := c.DBConfig(context.Background())
	require.NoError(t, err)

	assert.Equal(t, sd.Port, got.Port)
	assert.Equal(t, sd.Address, got.Host)
	assert.Equal(t, c.PostgresUser, got.User)
	assert.Equal(t, c.PostgresPass, got.Password)
	assert.Equal(t, c.PostgresDB, got.DBName)

	sd.Err = errors.New("negative test case")
	_, err = c.DBConfig(context.Background())
	assert.Error(t, err)
}

func TestConfigLoggerConfig(t *testing.T) {
	sd := &MockSD{
		Address: "moon",
		Port:    8888,
	}
	c := Config{
		LoggerPassSecret: "secretPass",
		LoggerPassSHA2:   "SHA2Pass",
		LoggerOutput:     "Stdout",
		LoggerLevel:      "info",
		LoggerType:       "async",
		sd:               sd,
	}

	got, err := c.LoggerConfig(context.Background())
	require.NoError(t, err)

	assert.Equal(t, sd.Port, got.Port)
	assert.Equal(t, sd.Address, got.Host)
	assert.Equal(t, c.LoggerPassSecret, got.PassSecret)
	assert.Equal(t, c.LoggerPassSHA2, got.PassSHA2)
	assert.Equal(t, c.LoggerOutput, got.Output)
	assert.Equal(t, c.LoggerLevel, got.Level)
	assert.Equal(t, c.LoggerType, got.Type)

	sd.Err = errors.New("negative test case")
	_, err = c.LoggerConfig(context.Background())
	assert.Error(t, err)
}

func TestConfigEmailConfig(t *testing.T) {
	c := Config{
		EmailAddress:  "user.namager@gmail.com",
		EmailPassword: "lvl484Golang",
		EmailHost:     "smtp.gmail.com",
		EmailPort:     587,
	}

	got, err := c.EmailConfig()
	require.NoError(t, err)

	assert.Equal(t, c.EmailAddress, got.Sender)
	assert.Equal(t, c.EmailPassword, got.Password)
	assert.Equal(t, c.EmailHost, got.Host)
	assert.Equal(t, c.EmailPort, got.Port)
}
