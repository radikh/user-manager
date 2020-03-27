// Package config is responsible for loading user-manager application config.

package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/lvl484/user-manager/logger"
	"github.com/lvl484/user-manager/config"
)

func TestNewConfig(t *testing.T) {
	type args struct {
		configName string
		configPath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"TestCorrectInput", args{configName: "viper.config", configPath: "./"}, false,
		}, {
			"TestIncorrectInput", args{configName: "not a viper", configPath: "my_path"}, true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewConfig(tt.args.configName, tt.args.configPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewViperCfg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNewLoggerConfig(t *testing.T) {
	v, err := NewConfig("testvipercfg", "./testdata/")
	if err != nil {
		t.Errorf("Cant start test, err: %v", err)
	}

	type fields struct {
		v *viper.Viper
	}

	tests := []struct {
		name   string
		fields fields
		want   *logger.LogConfig
	}{
		{
			name:   "test",
			fields: fields{v: v.v},
			want: &logger.LogConfig{
				Host:       "GREYLOGHOST",
				Port:       "77777",
				PassSecret: "secretpassword",
				PassSHA2:   "&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&",
				Output:     "Graylog",
				Level:      "info",
				Type:       "async",
			},
		}, {
			name:   "testFail",
			fields: fields{v: v.v},
			want: &logger.LogConfig{
				Host:       "GREYLOGHOST",
				Port:       "77777",
				PassSecret: "wrongpassword",
				PassSHA2:   "anotherwrongpassword",
				Output:     "Graylog",
				Level:      "info",
				Type:       "async",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := &Config{
				v: tt.fields.v,
			}
			got := conf.NewLoggerConfig()
			assert.Equal(t, tt.want, got)
		})

func TestNewConfig_Required(t *testing.T) {
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

	os.Unsetenv("POSTGRES_USER")

	cfg, err = config.NewConfig()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "POSTGRES_USER")
	assert.Nil(t, cfg)
}

func TestNewConfig_Default(t *testing.T) {
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
			got:      cfg.BindIP,
			expected: "0.0.0.0",
		}, {
			name:     "BIND_PORT",
			got:      cfg.BindPort,
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
