package storage

import (
	"testing"
)

func TestConnectToDB(t *testing.T) {
	conf := &DBConfig{
		Host:     "127.0.0.1",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		DBName:   "um_db",
	}

	incorrectConf := &DBConfig{
		Host:     "host",
		Port:     1111,
		User:     "user",
		Password: "pasword",
		DBName:   "DBName",
	}

	zeroConf := &DBConfig{}

	tests := []struct {
		name   string
		config *DBConfig
		expect bool
	}{
		{
			name:   "CorrectInput",
			config: conf,
			expect: true,
		},
		{
			name:   "IncorrectInput",
			config: incorrectConf,
			expect: false,
		},
		{
			name:   "ZeroInput",
			config: zeroConf,
			expect: false,
		},
	}
	for _, test := range tests {
		_, get := ConnectToDB(test.config)

		if (get != nil) == test.expect {
			t.Errorf("ConnectToDB(%v) expect %v", test.name, test.expect)
			return
		}
	}
}
