package config_test

import (
	"os"
	"testing"

	"github.com/lvl484/user-manager/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfigRequired(t *testing.T) {
	os.Setenv("POSTGRES_USER", "postgres")
	os.Setenv("POSTGRES_PASSWORD", "1q2w3e4r")
	os.Setenv("POSTGRES_DB", "um_db")
	os.Setenv("CONSUL_ADDRESS", "consul:8500")
	os.Setenv("CONSUL_TOKEN", "token")

	cfg, err := config.NewConfig()
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

	cfg, err := config.NewConfig()
	require.NoError(t, err)

	testCases := []struct {
		name     string
		got      interface{}
		expected interface{}
	}{
		{
			name:     "BIND_IP",
			got:      cfg.HTTPIP,
			expected: "0.0.0.0",
		}, {
			name:     "BIND_PORT",
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

	cfg, err := config.NewConfig()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "POSTGRES_USER")
	assert.Nil(t, cfg)
}
