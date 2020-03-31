// Package storage has to provide an interface to create a new postgres client
// with ability to reconnect after connection failure.
package storage

import (
<<<<<<< HEAD
	"database/sql"
	"fmt"

=======
>>>>>>> master
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
<<<<<<< HEAD

const pgStr = "host=%v port=%v user=%v password=%v dbname=%v sslmode=disable"

//DB connector
func ConnectToDB(pg *DBConfig) (*sql.DB, error) {
	pgConfig := fmt.Sprintf(pgStr, pg.Host, pg.Port, pg.User, pg.Password, pg.DBName)

	database, err := sql.Open("postgres", pgConfig)
	if err != nil {
		return nil, err
	}

	err = database.Ping()

	if err != nil {
		return nil, err
	}

	return database, nil
}
=======
>>>>>>> master
