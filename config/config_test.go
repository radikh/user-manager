package config_test

import (
	"os"
	"testing"

	"github.com/lvl484/user-manager/config"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig_Success(t *testing.T) {
	os.Setenv("POSTGRES_USER", "postgres")
	os.Setenv("POSTGRES_PASSWORD", "1q2w3e4r")
	os.Setenv("POSTGRES_DB", "um_db")

	os.Setenv("CONSUL_ADDRESS", "consul:8500")
	os.Setenv("CONSUL_TOKEN", "token")

	cfg, err := config.NewConfig()
	if err != nil {
		t.Fatal("want nil got error", err)
	}

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
		assert.Equal(t, tt.got, tt.expected, "The two values should be equal.")
	}
}

func TestNewConfig_Fail(t *testing.T) {
	os.Unsetenv("POSTGRES_USER")

	cfg, err := config.NewConfig()
	if err == nil {
		t.Error("want error got nil")
	}

	if cfg != nil {
		t.Errorf("want nil got %v", cfg)
	}
}
