// Package config is responsible for loading user-manager application config.

package config

import (
	"testing"

	"github.com/lvl484/user-manager/logger"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
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
	}
}
