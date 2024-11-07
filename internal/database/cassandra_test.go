package database

import (
	"context"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/cassandra"
	"github.com/testcontainers/testcontainers-go/wait"
)

func mustStartCassandraContainer() (func(context.Context) error, error) {
	cassandraContainer, err := cassandra.Run(
		context.Background(),
		"cassandra:4.1.3",
		testcontainers.WithAfterReadyCommand(
			testcontainers.NewRawCommand([]string{"cqlsh", "-e", "CREATE KEYSPACE test WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor': 1};"}),
		),
		testcontainers.WithWaitStrategy(
			wait.ForLog("Created default superuser role 'cassandra'").
				WithOccurrence(1).
				WithStartupTimeout(30*time.Second),
		),
	)

	if err != nil {
		return nil, err
	}
	cassandraKeyspace = "test"

	host, err := cassandraContainer.Host(context.Background())
	if err != nil {
		return cassandraContainer.Terminate, err
	}

	port, err := cassandraContainer.MappedPort(context.Background(), "9042/tcp")
	if err != nil {
		return cassandraContainer.Terminate, err
	}

	cassandraHosts = []string{host}
	cassandraPort = port.Port()

	return cassandraContainer.Terminate, err
}

func TestCassandraInstance(t *testing.T) {
	srv := Cassandra()
	if srv == nil {
		t.Fatal("New() returned nil")
	}
}

func TestCassandraHealth(t *testing.T) {
	srv := Cassandra()
	if srv == nil {
		t.Fatal("New() returned nil")
	}
	health := srv.Health()
	if health["status"] != "up" {
		t.Fatalf("expected status to be up, got %s", health["status"])
	}
}

func TestCassandraClose(t *testing.T) {
	srv := Cassandra()
	if srv == nil {
		t.Fatal("New() returned nil")
	}
	if err := srv.Close(); err != nil {
		t.Fatalf("Close() returned an error: %v", err)
	}
}
