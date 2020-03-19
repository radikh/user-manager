// Package config is responsible for loading user-manager application config.
// Basic configuration like consul credentials and address, http port to listen for requests,
// postgres schema name, credentials, and client timeout are read from environment variables.
package config

import (
	"strings"

	"github.com/lvl484/user-manager/logger"
	"github.com/spf13/viper"
)

// NewPostgresConfig returns pointer to PointerConfig with data read from viper.config.json
func NewLoggerConfig(configName, configPath string) (*logger.LogConfig, error) {
	v := viper.New()
	v.AddConfigPath(configPath)
	v.SetConfigName(configName)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	return &logger.LogConfig{
		Host:        v.GetString("logger.Host"),
		Port:        v.GetString("logger.Port"),
		Pass_Secret: v.GetString("logger.Pass_Secret"),
		Pass_SHA2:   v.GetString("logger.Pass_SHA2"),
		Output:      v.GetString("logger.Output"),
	}, nil
}
