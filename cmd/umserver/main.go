// Command umserver starts User Manager HTTP server.
// UM service stores user related context and credentials.
// It provides a REST API to perform a set of CRUD to manage users and an endpoint to authenticate.
// All users data will be stored in a database.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const GracefulShutdownTimeOut = 10 * time.Second

func main() {
	var (
		ctx, cancel = context.WithCancel(context.Background())
		wg          = new(sync.WaitGroup)
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
	err := gracefulShutdown(GracefulShutdownTimeOut, wg, srv)
	if err != nil {
		log.Fatalf("Server graceful shutdown failed: %v", err)
	}

	log.Println("Server was gracefully stopped!")
	os.Exit(code)
}

func gracefulShutdown(timeout time.Duration, wg *sync.WaitGroup, srv *http.Server) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	// Shutdown HTTP server
	go func() {
		err := srv.Shutdown(ctx)
		if err != nil {
			log.Println("shutdown error: %w", err)
		}
	}()

	// Wait with timeout
	go func() {
		defer cancel()
		wg.Wait()
	}()

	<-ctx.Done()

	if ctx.Err() == context.DeadlineExceeded {
		return errors.New("timeout")
	}

	return nil
}
