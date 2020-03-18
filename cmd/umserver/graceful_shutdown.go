package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

const GracefulShutdownTimeOut = 10 * time.Second

func GracefulShutdown(ctx context.Context, cancel context.CancelFunc, interrupt chan os.Signal, srv *http.Server) {
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
}

// WaitTimeout waits for the waitgroup for the specified max timeout.
// Returns information about waiting group.
func WaitTimeout(ctx context.Context, wg *sync.WaitGroup) {
	c := make(chan struct{})

	go func() {
		defer close(c)
		wg.Wait()
	}()

	select {
	case <-c:
		log.Println("Wait group finished")
		return
	case <-ctx.Done():
		log.Println("Timed out waiting for wait group")
		return
	}
}
