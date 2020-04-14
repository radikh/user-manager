package main

import (
	"context"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/lvl484/user-manager/logger"
	"github.com/lvl484/user-manager/server"
)

func gracefulShutdown(timeout time.Duration, wg *sync.WaitGroup, srv *server.HTTP, closers ...io.Closer) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	// Shutdown HTTP server
	go func() {
		err := srv.Stop(ctx)
		if err != nil {
			logger.LogUM.Errorf("shutdown error: %w", err)
		}

		for _, component := range closers {
			err := component.Close()
			if err != nil {
				logger.LogUM.Errorf("component error: %w", err)
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
