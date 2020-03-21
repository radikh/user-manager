package main

import (
	"database/sql"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/connect"
	"github.com/segmentio/kafka-go"
)

func TestGracefulShutdown_Success(t *testing.T) {
	relatedComponents := setComponents()

	err := gracefulShutdown(3*time.Second, new(sync.WaitGroup), new(http.Server), relatedComponents)
	if err != nil {
		t.Error(err)
	}
}

func TestGracefulShutdown_Fail(t *testing.T) {
	relatedComponents := setComponents()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	err := gracefulShutdown(1*time.Nanosecond, wg, new(http.Server), relatedComponents)
	if err == nil {
		t.Error("want error got nil")
	}
}

func Test_GracefulShutdown_Timeout(t *testing.T) {
	relatedComponents := setComponents()

	t.Run("with timeout", func(t *testing.T) {
		err := gracefulShutdown(0*time.Second, new(sync.WaitGroup), new(http.Server), relatedComponents)
		if err.Error() != "timeout" {
			t.Fatal("want timeout got nil")
		}
	})
}

// setComponents is using for setting database, consul and kafka as related component
func setComponents() *closingStructure {
	// There will be actual information about PostgreSQL connection in future
	connStr := "user=postgres password=postgres dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	if err = db.Close(); err != nil {
		log.Fatal(err)
	}

	// There will be actual information in future
	netConnection, err := net.Dial("tcp", "golang.org:80")
	if err != nil {
		log.Fatalf("net connection failed:%+v", err)
	}

	// There will be actual information about kafka in future
	kafkaConnection := kafka.NewConn(netConnection, "topic_name", 0)

	// There will be actual information about consul in future
	consulClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatalf("consul client failed: %+v", err)
	}

	consulService, err := connect.NewService("service_name", consulClient)
	if err != nil {
		log.Fatalf("consul service failed:%+v", err)
	}

	// relatedComponents structure includes all components, which are connected with our application
	// The following components: db connection, kafkaConnection and consulService
	relatedComponents := closingStructure{
		components: []io.Closer{
			db,
			kafkaConnection,
			consulService,
		},
	}
	return &relatedComponents
}
