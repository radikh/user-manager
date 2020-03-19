// Package logger provides implementetion of writing Log messages to Graylog, file, and stdout.
// Supports log levels and destination.
package logger

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	loggerHost        = "logger_um.Host"
	loggerPort        = "logger_um.Port"
	loggerPass_Secret = "logger_um.Pass_Secret"
	loggerPass_SHA2   = "logger_um.Pass_SHA2"
	loggerOutput      = "logger_um.Output"
)

func TestNullFormatter_Format(t *testing.T) {
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NullFormatter.Format() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewLogger(t *testing.T) {
	v := viper.New()
	v.AddConfigPath("../config/")
	v.SetConfigName("viper.config")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		t.Error(err)
	}

	_, err := os.OpenFile("user_manager_api.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}

	conf := &LogConfig{
		Host:        v.GetString(loggerHost),
		Port:        v.GetString(loggerPort),
		Pass_Secret: v.GetString(loggerPass_Secret),
		Pass_SHA2:   v.GetString(loggerPass_SHA2),
		Output:      v.GetString(loggerOutput),
	}
	incorrectConf := &LogConfig{
		Host:        "locallviv",
		Port:        "15000",
		Pass_Secret: "root",
		Pass_SHA2:   "asdfsdfdsfewffsdvsvdsvfdsvsvsd",
		Output:      "Greenlog",
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
			err := NewLogger(tt.lc)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLogger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMessage(t *testing.T) {
	logger := log.New()
	logger.Out = &bytes.Buffer{}
	_ = log.NewEntry(logger)
	log.SetLevel(log.PanicLevel)

	type args struct {
		level log.Level
		m     string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Debug",
			args: args{level: log.DebugLevel, m: "Debug level"},
		}, {
			name: "Warn",
			args: args{level: log.WarnLevel, m: "Warn level"},
		}, {
			name: "Error",
			args: args{level: log.ErrorLevel, m: "Error level"},
		}, {
			name: "Info",
			args: args{level: log.InfoLevel, m: "Info level"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Message(tt.args.level, tt.args.m)
		})
	}
}
