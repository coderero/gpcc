package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *CassandraServer) RegisterRoutes() http.Handler {
	r := gin.Default()

	return r
}
