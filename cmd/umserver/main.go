// Command umserver starts User Manager HTTP server.
// UM service stores user related context and credentials.
// It provides a REST API to perform a set of CRUD to manage users and an endpoint to authenticate.
// All users data will be stored in a database.
package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/lvl484/user-manager/config"
	"github.com/lvl484/user-manager/logger"
	"github.com/lvl484/user-manager/storage"
	"github.com/sirupsen/logrus"
)

const gracefulShutdownTimeOut = 10 * time.Second

var configuration *config.Config
var Log logger.Logger

func init() {
	var err error
	configuration, err = config.NewConfig("viper.config", "../../config")
	if err != nil {
		fmt.Println(err)
	}
	loggerConfig := configuration.NewLoggerConfig()
	logUM := logrus.New()
	err = logger.ConfigLogger(logUM, loggerConfig)
	if err != nil {
		fmt.Println(err)
	}
	logger.SetLogger(logUM)
	Log = logger.LogUM
}

func main() {
	var (
		ctx, cancel = context.WithCancel(context.Background())
		wg          = new(sync.WaitGroup)
		closers     []io.Closer
	)
	var err error

	// TODO: Replace with HTTP server implemented in server package
	srv := &http.Server{Addr: ":8099"}
	Log.Info("Server start at %s...", srv.Addr)
	fmt.Println("A am started")
	// Go routine with run HTTP server
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer cancel()

		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			Log.Error("%v\n", err)
		}
	}()
	Log.Info("Server Listening at %s...", srv.Addr)

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
		Log.Error("%v\n", err)
	}
	Log.Info("Successfully connected to %s", storage.DBConfig.DBName)
	defer db.Close()

	// TODO: There will be actual information about consul in future
	// ...
	// TODO: There will be actual information about kafka in future

	// Watch errors and os signals
	interrupt, code := make(chan os.Signal, 1), 0
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-interrupt:
		Log.Info("Pressed Ctrl+C to terminate server...")
		cancel()
	case <-ctx.Done():
		code = 1
	}

	Log.Info("Server is Stopping...")

	// Stop application
	err = gracefulShutdown(gracefulShutdownTimeOut, wg, srv, closers...)
	if err != nil {
		Log.Panicf("Server graceful shutdown failed: %v", err)
	}

	Log.Info("Server was gracefully stopped!")
	os.Exit(code)
}
