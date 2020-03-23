// Package config is responsible for loading user-manager application config.
// Basic configuration like consul credentials and address, http port to listen for requests,
// postgres schema name, credentials, and client timeout are read from environment variables.
package config

import (
	"reflect"
	"testing"

	"github.com/lvl484/user-manager/logger"

	"github.com/spf13/viper"
)

func TestNewViperCfg(t *testing.T) {
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
			_, err := NewViperCfg(tt.args.configName, tt.args.configPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewViperCfg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestViperCfg_NewLoggerConfig(t *testing.T) {
	v, err := NewViperCfg("testvipercfg", "./testdata/")
	if err != nil {
		t.Errorf("Cant start test, err: %v", err)
	}

	type fields struct {
		v *viper.Viper
	}

	tests := []struct {
		name    string
		fields  fields
		want    *logger.LogConfig
		wantErr bool
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
			wantErr: false,
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
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vpcfg := &ViperCfg{
				v: tt.fields.v,
			}
			got, err := vpcfg.NewLoggerConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("ViperCfg.NewLoggerConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ViperCfg.NewLoggerConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
