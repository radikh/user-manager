// Package storage has to provide an interface to create a new postgres client
// with ability to reconnect after connection failure.
package storage

import (
	"database/sql"
	"fmt"

	// _ is used for registering the pq driver as a database driver,
	// without importing any other functions
	_ "github.com/lib/pq"
)

// Config of Postgres DB
type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

const pgStr = "host=%v port=%v user=%v password=%v dbname=%v sslmode=disable"
const DBDriverName = "postgres"

//DB connector
func ConnectToDB(pg *DBConfig) (*sql.DB, error) {
	pgConfig := GetDBConfigString(pg)
	database, err := sql.Open(DBDriverName, pgConfig)

	if err != nil {
		return nil, err
	}

	err = database.Ping()

	if err != nil {
		return nil, err
	}

	return database, nil
}

func GetDBConfigString(pg *DBConfig) string {
	return fmt.Sprintf(pgStr, pg.Host, pg.Port, pg.User, pg.Password, pg.DBName)
}
