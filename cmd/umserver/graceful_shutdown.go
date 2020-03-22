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

func gracefulShutdown(timeout time.Duration, wg *sync.WaitGroup, srv *http.Server, closers ...io.Closer) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	// Shutdown HTTP server
	go func() {
		err := srv.Shutdown(ctx)
		if err != nil {
			log.Println("shutdown error: %w", err)
		}

		for _, component := range closers {
			err := component.Close()
			if err != nil {
				log.Println("component error: %w", err)
			}
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
