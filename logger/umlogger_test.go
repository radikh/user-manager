// Package logger provides implementetion of writing Log messages to Graylog, file, and stdout.
// Supports log levels and destination.
package logger

import "testing"

func TestNewLogger(t *testing.T) {
	type args struct {
		lc *LogConfig
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewLogger(tt.args.lc)
		})
	}
}
