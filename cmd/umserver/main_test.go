package main

import (
	"errors"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"
)

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

	err := gracefulShutdown(3*time.Second, new(sync.WaitGroup), new(http.Server), closers...)
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
	err := gracefulShutdown(1*time.Nanosecond, wg, new(http.Server), closers...)

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
			return errors.New("can't close connection")
		}),
		NewCloserMock("kafka", func() error {
			return nil
		}),
	}

	t.Run("with timeout", func(t *testing.T) {
		err := gracefulShutdown(0*time.Second, new(sync.WaitGroup), new(http.Server), closers...)
		if err.Error() != "timeout" {
			t.Fatal("want timeout got nil")
		}
	})
}
