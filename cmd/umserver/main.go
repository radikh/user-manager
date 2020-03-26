// Command umserver starts User Manager HTTP server.
// UM service stores user related context and credentials.
// It provides a REST API to perform a set of CRUD to manage users and an endpoint to authenticate.
// All users data will be stored in a database.
package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/lvl484/user-manager/storage"
)

const gracefulShutdownTimeOut = 10 * time.Second

func main() {
	var (
		ctx, cancel = context.WithCancel(context.Background())
		wg          = new(sync.WaitGroup)
		closers     []io.Closer
	)

	// TODO: Replace with HTTP server implemented in server package
	srv := &http.Server{Addr: ":8099"}

	// Go routine with run HTTP server
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer cancel()

		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Printf("%v\n", err)
		}
	}()
	log.Printf("Server Listening at %s...", srv.Addr)

	// pgConfig will be taken from package config, but it hasn't ready yet
	pgConfig := storage.DBConfig{
		Host:     "127.0.0.1",
		Port:     "5432",
		User:     "postgres",
		Password: "postgres",
		DBName:   "um_db",
	}

	db, err := storage.ConnectToDB(&pgConfig)
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	// TODO: There will be actual information about consul in future
	// ...
	// TODO: There will be actual information about kafka in future

	// Watch errors and os signals
	interrupt, code := make(chan os.Signal, 1), 0
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-interrupt:
		log.Print("Pressed Ctrl+C to terminate server...")
		cancel()
	case <-ctx.Done():
		code = 1
	}

	log.Print("Server is Stopping...")

	// Stop application
	err = gracefulShutdown(gracefulShutdownTimeOut, wg, srv, closers...)
	if err != nil {
		log.Fatalf("Server graceful shutdown failed: %v", err)
	}

	log.Println("Server was gracefully stopped!")
	os.Exit(code)
}
