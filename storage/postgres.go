// Package storage has to provide an interface to create a new postgres client
// with ability to reconnect after connection failure.
package storage

import (
	"database/sql"
	"fmt"
	"log"

	// _ mean that we can use all function from this package
	_ "github.com/lib/pq"
)

// Config of Postgres DB
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

const pgStr = "host=%v port=%v user=%v password=%v dbname=%v sslmode=disable"

//DB connector
func ConnectToDB(pg *DBConfig) (*sql.DB, error) {
	pgConfig := fmt.Sprintf(pgStr, pg.Host, pg.Port, pg.User, pg.Password, pg.DBName)

	database, err := sql.Open("postgres", pgConfig)
	if err != nil {
		log.Print("Could not connect to ", pg.DBName)
		return nil, err
	}

	err = database.Ping()

	if err != nil {
		log.Print("Could not connect to ", pg.DBName)
		return nil, err
	}

	log.Print("Successfully connected to ", pg.DBName)

	return database, nil
}
