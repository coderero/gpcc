package server

import (
	"connector/internals/database"
	"fmt"
	"net/http"
	"time"
)

// The `CassandraServer` struct type represents a server for a Cassandra database.
type CassandraServer struct {
	port int

	db database.Service
}

// The function `NewCassandraServer` creates a new HTTP server for a Cassandra database with specified
// port and timeouts.
func NewCassandraServer(port int) *http.Server {
	NewServer := &CassandraServer{
		port: port,

		db: database.Cassandra(),
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
