package storage

import (
	"testing"

	_ "github.com/lib/pq"
)

func TestConnectToDB(t *testing.T) {
	conf := &PgClient{
		Host:     "127.0.0.1",
		Port:     "5432",
		User:     "postgres",
		Password: "postgres",
		DBName:   "um_db",
	}

	incorrectConf := &PgClient{
		Host:     "host",
		Port:     "port",
		User:     "user",
		Password: "pasword",
		DBName:   "DBName",
	}

	zeroConf := &PgClient{}

	tests := []struct {
		name   string
		config *PgClient
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
