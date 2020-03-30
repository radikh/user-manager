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

	"github.com/lvl484/user-manager/config"
	"github.com/lvl484/user-manager/logger"

	_ "github.com/lib/pq"
)

const gracefulShutdownTimeOut = 10 * time.Second

func main() {
	var (
		ctx, cancel = context.WithCancel(context.Background())
		wg          = new(sync.WaitGroup)
		closers     []io.Closer
	)

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	loggerConfig, err := cfg.LoggerConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = logger.SetLogger(loggerConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Example
	log.Println(cfg)
	log.Println(cfg.LoggerConfig(ctx))
	log.Println(cfg.DBConfig(ctx))

	// TODO: Replace with HTTP server implemented in server package
	srv := &http.Server{
		Addr:         cfg.ServerAddress(),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	// Go routine with run HTTP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()

		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.LogUM.Error("%v\n", err)
		}
	}()
	logger.LogUM.Info("Server Listening at %s...", srv.Addr)

	// TODO: There will be actual information about PostgreSQL connection in future
	// ...
	// TODO: There will be actual information about consul in future
	// ...
	// TODO: There will be actual information about kafka in future

	// Watch errors and os signals
	interrupt, code := make(chan os.Signal, 1), 0
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-interrupt:
		logger.LogUM.Info("Pressed Ctrl+C to terminate server...")
		cancel()
	case <-ctx.Done():
		code = 1
	}

	logger.LogUM.Info("Server is Stopping...")

	// Stop application
	err = gracefulShutdown(gracefulShutdownTimeOut, wg, srv, closers...)
	if err != nil {
		logger.LogUM.Fatalf("Server graceful shutdown failed: %v", err)
	}

	logger.LogUM.Info("Server was gracefully stopped!")
	os.Exit(code)
}
