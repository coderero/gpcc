package database

import (
	"context"
	"log"
	"os"
	"testing"
)

// The TestMain function sets up and tears down containers for PostgreSQL and Cassandra before and
// after running tests.
func TestMain(m *testing.M) {
	postgresTeardown, err := mustStartPostgresContainer()
	if err != nil {
		log.Fatalf("could not start postgres container: %v", err)
	}

	cassandraTeardown, err := mustStartCassandraContainer()
	if err != nil {
		log.Fatalf("could not start cassandra container: %v", err)
	}

	code := m.Run()

	if err := postgresTeardown(context.Background()); err != nil {
		log.Fatalf("could not stop postgres container: %v", err)
	}

	if err := cassandraTeardown(context.Background()); err != nil {
		log.Fatalf("could not stop cassandra container: %v", err)
	}

	os.Exit(code)
}
