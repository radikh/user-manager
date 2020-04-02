package storage

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const dbUser = "POSTGRES_USER"
const dbPassword = "POSTGRES_PASSWORD"
const db = "POSTGRES_DB"
const expectedLine = "host=localhost port=5432 user=postgres password=postgres dbname=um_db sslmode=disable"

func TestConnectToDB(t *testing.T) {
	user, _ := os.LookupEnv(dbUser)
	password, _ := os.LookupEnv(dbPassword)
	dbname, _ := os.LookupEnv(db)

	conf := &DBConfig{
		Host:     "localhost",
		Port:     5432,
		User:     user,
		Password: password,
		DBName:   dbname,
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
		assert.Equal(t, get == nil, test.expect)
	}
}

func TestGetDBConfigString(t *testing.T) {
	conf := &DBConfig{
		Host:     "localhost",
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
		get := getDBConfigString(test.config)
		assert.Equal(t, get == expectedLine, test.expect)
	}
}
