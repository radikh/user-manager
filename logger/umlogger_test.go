// Package logger provides implementetion of writing Log messages to Graylog, file, and stdout.
// Supports log levels and destination.
package logger

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const (
	loggerHost       = "loggerUM.Host"
	loggerPort       = "loggerUM.Port"
	loggerPassSecret = "loggerUM.PassSecret"
	loggerPassSHA2   = "loggerUM.PassSHA2"
	loggerOutput     = "loggerUM.Output"
)

func TestNullFormatterFormat(t *testing.T) {
	type args struct {
		e *log.Entry
	}
	formatter := &log.JSONFormatter{}
	b, err := formatter.Format(log.WithField("opg", errors.New("user managment test")))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry1 := make(map[string]interface{})
	err = json.Unmarshal(b, &entry1)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}
	logger := log.New()
	logger.Out = &bytes.Buffer{}
	entry := log.NewEntry(logger)

	tests := []struct {
		name    string
		n       log.Formatter
		args    args
		want    []byte
		wantErr bool
	}{
		{name: "Null format",
			n:       NullFormatter{},
			args:    args{e: entry},
			want:    []byte{},
			wantErr: false,
		}, {
			name:    "JSON format",
			n:       NullFormatter{},
			args:    args{e: &log.Entry{Logger: logger, Data: entry1}},
			want:    []byte{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NullFormatter{}
			got, err := n.Format(tt.args.e)
			if (err != nil) != tt.wantErr {
				t.Errorf("NullFormatter.Format() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConfigLogger(t *testing.T) {
	v := viper.New()
	v.AddConfigPath("../config/testdata/")
	v.SetConfigName("testvipercfg")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		t.Error(err)
	}
	logger := log.New()
	conf := &LogConfig{
		Host:       v.GetString(loggerHost),
		Port:       v.GetString(loggerPort),
		PassSecret: v.GetString(loggerPassSecret),
		PassSHA2:   v.GetString(loggerPassSHA2),
		Output:     v.GetString(loggerOutput),
	}
	incorrectConf := &LogConfig{
		Host:       "locallviv",
		Port:       "15000",
		PassSecret: "root",
		PassSHA2:   "asdfsdfdsfewffsdvsvdsvfdsvsvsd",
		Output:     "Greenlog",
	}
	conf_file := &LogConfig{
		Output: "File",
	}
	conf_stdout := &LogConfig{
		Output: "Stdout",
	}

	tests := []struct {
		name    string
		lc      *LogConfig
		wantErr bool
	}{
		{
			name:    "CorrectConfig1",
			lc:      conf,
			wantErr: false,
		}, {
			name:    "CorrectConfig2",
			lc:      conf_file,
			wantErr: false,
		}, {
			name:    "CorrectConfig3",
			lc:      conf_stdout,
			wantErr: false,
		}, {
			name:    "UncorrectConfig",
			lc:      incorrectConf,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ConfigLogger(logger, tt.lc)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLogger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestLogConfigSetLoggerToFile(t *testing.T) {

	_, err := os.OpenFile("user_manager_api.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		t.Errorf("LogConfig.setLoggerToFile() error = %v", err)
	}
}

func TestLogConfigSetLoggerToStdout(t *testing.T) {
	logger := log.New()
	conf_file := &LogConfig{
		Output: "Filename",
	}
	conf_stdout := &LogConfig{
		Output: "Stdout",
	}

	conf_stdout.setLoggerToStdout(logger)
	assert.Equal(t, os.Stdout, log.StandardLogger().Out)
	conf_file.setLoggerToStdout(logger)
	assert.NotEqual(t, os.Stdout, log.StandardLogger().Out)

}

func TestLogConfigSetLoggerToGraylog(t *testing.T) {
	v := viper.New()
	v.AddConfigPath("../config/testdata/")
	v.SetConfigName("testvipercfg")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		t.Error(err)
	}
	logger := log.New()
	conf := &LogConfig{
		Host:       v.GetString(loggerHost),
		Port:       v.GetString(loggerPort),
		PassSecret: v.GetString(loggerPassSecret),
		PassSHA2:   v.GetString(loggerPassSHA2),
		Output:     v.GetString(loggerOutput),
	}
	incorrectConf := &LogConfig{
		Host:       "locallviv",
		Port:       "15000",
		PassSecret: "root",
		PassSHA2:   "asdfsdfdsfewffsdvsvdsvfdsvsvsd",
		Output:     "Greenlog",
	}

	tests := []struct {
		name    string
		lc      *LogConfig
		wantErr bool
	}{
		{
			name:    "CorrectConfig1",
			lc:      conf,
			wantErr: false,
		}, {
			name:    "IncorrectConfig",
			lc:      incorrectConf,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ConfigLogger(logger, tt.lc)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLogger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
