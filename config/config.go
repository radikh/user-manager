// Package config is responsible for loading user-manager application config.
// Basic configuration like consul credentials and address, http port to listen for requests,
// postgres schema name, credentials, and client timeout are read from environment variables.
package config

import (
	"strings"

	"github.com/lvl484/user-manager/logger"

	"github.com/spf13/viper"
)

const (
	loggerHost       = "loggerUM.Host"
	loggerPort       = "loggerUM.Port"
	loggerPassSecret = "loggerUM.PassSecret"
	loggerPassSHA2   = "loggerUM.PassSHA2"
	loggerOutput     = "loggerUM.Output"
)

type ViperCfg struct {
	v *viper.Viper
}

func NewViperCfg(configName, configPath string) (*ViperCfg, error) {
	v := viper.New()
	v.AddConfigPath(configPath)
	v.SetConfigName(configName)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	return &ViperCfg{v: v}, nil
}

// NewPostgresConfig returns pointer to PointerConfig with data read from viper.config.json
func (vpcfg *ViperCfg) NewLoggerConfig() (*logger.LogConfig, error) {
	return &logger.LogConfig{
		Host:       vpcfg.v.GetString(loggerHost),
		Port:       vpcfg.v.GetString(loggerPort),
		PassSecret: vpcfg.v.GetString(loggerPassSecret),
		PassSHA2:   vpcfg.v.GetString(loggerPassSHA2),
		Output:     vpcfg.v.GetString(loggerOutput),
	}, nil
}
