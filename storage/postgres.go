// Package storage has to provide an interface to create a new postgres client
// with ability to reconnect after connection failure.
package storage

import (
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
