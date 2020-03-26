package config_test

import (
	"os"
	"testing"

	"github.com/lvl484/user-manager/config"
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

	if cfg.PostgresUser != "postgres" {
		t.Errorf("want postgres got %s", cfg.PostgresUser)
	}
	if cfg.PostgresPass != "1q2w3e4r" {
		t.Errorf("want 1q2w3e4r got %s", cfg.PostgresPass)
	}
	if cfg.PostgresDB != "um_db" {
		t.Errorf("want um_db got %s", cfg.PostgresDB)
	}
	if cfg.ConsulAddress != "consul:8500" {
		t.Errorf("want consul:8500 got %s", cfg.ConsulAddress)
	}
	if cfg.ConsulToken != "token" {
		t.Errorf("want token got %s", cfg.ConsulToken)
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
