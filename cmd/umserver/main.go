// Command umserver starts User Manager HTTP server.
// UM service stores user related context and credentials.
// It provides a REST API to perform a set of CRUD to manage users and an endpoint to authenticate.
// All users data will be stored in a database.
package main

import (
	"context"
	"fmt"
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
	log.Print("Server Listening...")

	select {
	case <-interrupt:
		log.Print("Pressed Ctrl+C to terminate server...")
		cancel()
	case <-ctx.Done():
		log.Print(ctx.Err())
	}

	log.Print("Server is Stopping...")

	ctx, cancel = context.WithTimeout(context.Background(), GracefulShutdownTimeOut)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if waitTimeout(ctx, wg) {
		fmt.Println("Timed out waiting for wait group")
	} else {
		fmt.Println("Wait group finished")
	}
	log.Println("Server was successful stopped!")
}

// waitTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out.
func waitTimeout(ctx context.Context, wg *sync.WaitGroup) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-ctx.Done():
		return true // timed out
	}
}
