package main

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

// closingStructure includes all components, which are connected with our application
type closingStructure struct {
	components []io.Closer
}

func gracefulShutdown(timeout time.Duration, wg *sync.WaitGroup, srv *http.Server, closers *closingStructure) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	// Shutdown HTTP server
	go func() {
		for _, component := range closers.components {
			err := component.Close()
			if err != nil {
				log.Println("component error: %w", err)
			}
		}

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
