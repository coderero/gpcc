package server

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"connector/internals/database"
)

type PostgresServer struct {
	port int

	db database.Service
}

func NewPostgresServer(port int) *http.Server {
	NewServer := &PostgresServer{
		port: port,

		db: database.Postgres(),
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
