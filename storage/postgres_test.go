package storage

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	dbUser       = "POSTGRES_USER"
	dbPassword   = "POSTGRES_PASSWORD"
	db           = "POSTGRES_DB"
	expectedLineOk   = "host=localhost port=5432 user=postgres password=postgres dbname=um_db sslmode=disable"
	expectedLineZero = "host= port=0 user= password= dbname= sslmode=disable"
)

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
	
	zeroConf := &DBConfig{}

	bd, err := ConnectToDB(conf)
	assert.NotNil(t, bd)
	assert.Nil(t, err)
	bd, err = ConnectToDB(zeroConf)
	assert.Nil(t, bd)
	assert.NotNil(t, err)
}

func TestGetDBConfigString(t *testing.T) {
	conf := &DBConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		DBName:   "um_db",
	}

	zeroConf := &DBConfig{}

	tests := []struct {
		name   string
		config *DBConfig
		expect string
	}{
		{
			name:   "CorrectInput",
			config: conf,
			expect: expectedLineOk,
		},
		{
			name:   "ZeroInput",
			config: zeroConf,
			expect: expectedLineZero,
		},
	}
	for _, test := range tests {
		get := getDBConfigString(test.config)
		assert.Equal(t, get, test.expect)
	}
}
