// Package storage has to provide an interface to create a new postgres client
// with ability to reconnect after connection failure.
package storage

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// Config of Postgres DB
type PgClient struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
	Db       *sql.DB
}

const pgStr = "host=%v port=%v user=%v password=%v dbname=%v sslmode=disable"

//DB connector
func ConnectToDb(pg *PgClient) *sql.DB {

	pgConfig := fmt.Sprintf(pgStr, pg.Host, pg.Port, pg.User, pg.Password, pg.DbName)

	database, err := sql.Open("postgres", pgConfig)
	if err != nil {
		log.Print(err)
		log.Print("Could not connect to ", pg.DbName)
	}
	if pg.isAlive() != nil {
		log.Print(err)
		log.Print("Could not connect to ", pg.DbName)
	}
	log.Print("Successfully connected to ", pg.DbName)
	pg.Db = database

	return pg.Db
}

//Check connect to DB return err if no connection
func (pg *PgClient) isAlive() error {
	err := pg.Db.Ping()
	return err
}
