package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *PostgresServer) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.GET("/", s.HelloWorldHandler)

	r.GET("/health", s.healthHandler)

	return r
}

func (s *PostgresServer) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *PostgresServer) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
