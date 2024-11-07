package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocql/gocql"
	_ "github.com/joho/godotenv/autoload"
)

type cassandraService struct {
	session *gocql.Session
}

var (
	cassandraKeyspace = os.Getenv("CASSANDRA_KEYSPACE")
	cassandraUsername = os.Getenv("CASSANDRA_USERNAME")
	cassandraPassword = os.Getenv("CASSANDRA_PASSWORD")
	cassandraHosts    = strings.Split(os.Getenv("CASSANDRA_HOSTS"), ",")
	cassandraPort     = os.Getenv("CASSANDRA_PORT")
	cassandraInstance *cassandraService
)

// Cassandra returns a new Cassandra service.
func Cassandra() Service {
	// Reuse Connection
	if cassandraInstance != nil {
		return cassandraInstance
	}

	// Create a new Cassandra cluster
	cluster := gocql.NewCluster(cassandraHosts...)
	cluster.Keyspace = cassandraKeyspace
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: cassandraUsername,
		Password: cassandraPassword,
	}

	if cassandraPort != "" {
		port, err := strconv.Atoi(cassandraPort)
		if err != nil {
			log.Fatalf("could not convert cassandra port to int: %v", err)
		}
		cluster.Port = port
	} else {
		cluster.Port = 9042
	}
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	cassandraInstance = &cassandraService{
		session: session,
	}
	return cassandraInstance
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *cassandraService) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Try to execute a simple query to verify connection
	err := s.session.Query("SELECT now() FROM system.local").WithContext(ctx).Exec()
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Printf("Cassandra is down: %v", err) // Log the error without terminating
		return stats
	}

	// Database is up
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Add additional checks for node metrics (e.g., number of nodes)
	iter := s.session.Query("SELECT count(*) FROM system.peers").WithContext(ctx).Iter()
	var nodeCount int
	if iter.Scan(&nodeCount) {
		stats["node_count"] = fmt.Sprintf("%d", nodeCount+1) // Including local node
	} else {
		stats["node_count"] = "unknown"
		stats["warning"] = "Could not retrieve node count"
	}
	iter.Close()

	// Add cluster name for context
	var clusterName string
	if err := s.session.Query("SELECT cluster_name FROM system.local").WithContext(ctx).Scan(&clusterName); err == nil {
		stats["cluster_name"] = clusterName
	} else {
		stats["cluster_name"] = "unknown"
		stats["warning"] = "Could not retrieve cluster name"
	}

	return stats
}

// Close terminates the database connection.
func (s *cassandraService) Close() error {
	s.session.Close()
	return nil
}

// Instance returns the underlying database connection.
func (s *cassandraService) Instance() interface{} {
	return s.session
}
