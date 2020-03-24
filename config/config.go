// Package config is responsible for loading user-manager application config.

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
	LoggerLevel      = "loggerUM.Level"
	LoggerType       = "loggerUM.Type"
)

type Config struct {
	v *viper.Viper
}

func NewConfig(configName, configPath string) (*Config, error) {
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
