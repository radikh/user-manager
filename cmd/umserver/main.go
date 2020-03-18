// Command umserver starts User Manager HTTP server.
// UM service stores user related context and credentials.
// It provides a REST API to perform a set of CRUD to manage users and an endpoint to authenticate.
// All users data will be stored in a database.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	var (
		ctx, cancel = context.WithCancel(context.Background())
		interrupt   = make(chan os.Signal, 1)
		wg          = new(sync.WaitGroup)
	)

	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	// TODO: Replace with HTTP server implemented in server package
	srv := &http.Server{Addr: ":8099"}
	// Go routine with run Server
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()

		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Printf("%v\n", err)
		}
	}()
	log.Printf("Server Listening at%s...", srv.Addr)

	go func() {
		GracefulShutdown(ctx, cancel, interrupt, srv)
	}()

	WaitTimeout(ctx, wg)
	log.Println("Server was successful stopped!")
}
