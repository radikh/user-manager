package main

import (
	"errors"
	"io"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/lvl484/user-manager/config"
	"github.com/lvl484/user-manager/logger"
	"github.com/lvl484/user-manager/server"
)

func TestMain(m *testing.M) {
	logger.SetLogger(&logger.LogConfig{Output: "Stdout", Level: "debug"})

	code := m.Run()

	os.Exit(code)
}

func TestGracefulShutdown_Success(t *testing.T) {
	closers := []io.Closer{
		NewCloserMock("consul", func() error {
			return nil
		}),
		NewCloserMock("postgres", func() error {
			return errors.New("can't close connection")
		}),
		NewCloserMock("kafka", func() error {
			return nil
		}),
	}

	err := gracefulShutdown(3*time.Second, new(sync.WaitGroup), NewServerMock(), closers...)
	if err != nil {
		t.Error(err)
	}
}

func TestGracefulShutdown_Fail(t *testing.T) {
	closers := []io.Closer{
		NewCloserMock("consul", func() error {
			return nil
		}),
		NewCloserMock("postgres", func() error {
			return errors.New("can't close connection")
		}),
		NewCloserMock("kafka", func() error {
			return nil
		}),
	}

	wg := new(sync.WaitGroup)
	wg.Add(1)
	err := gracefulShutdown(1*time.Nanosecond, wg, NewServerMock(), closers...)
	if err == nil {
		t.Error("want error got nil")
	}
}

func Test_GracefulShutdown_Timeout(t *testing.T) {
	closers := []io.Closer{
		NewCloserMock("consul", func() error {
			return nil
		}),
		NewCloserMock("postgres", func() error {
			return nil
		}),
		NewCloserMock("kafka", func() error {
			return nil
		}),
	}

	t.Run("with timeout", func(t *testing.T) {
		err := gracefulShutdown(0*time.Second, new(sync.WaitGroup), NewServerMock(), closers...)
		if err.Error() != "timeout" {
			t.Fatal("want timeout got nil")
		}
	})
}

func NewServerMock() *server.HTTP {
	return server.NewHTTP(&config.Config{}, nil)
}

type CloserMock struct {
	name     string
	expected func() error
}

func NewCloserMock(name string, expected func() error) io.Closer {
	return CloserMock{
		name:     name,
		expected: expected,
	}
}

func (cm CloserMock) Close() error {
	return cm.expected()
}
