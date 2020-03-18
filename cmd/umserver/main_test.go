package main

import (
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestGracefulShutdown_Success(t *testing.T) {
	err := gracefulShutdown(3*time.Second, new(sync.WaitGroup), new(http.Server))
	if err != nil {
		t.Error(err)
	}
}

func TestGracefulShutdown_Fail(t *testing.T) {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	err := gracefulShutdown(1*time.Nanosecond, wg, new(http.Server))
	if err == nil {
		t.Error("want error got nil")
	}
}
