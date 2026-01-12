package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	Port string
	rdb  *redis.Client
}

func NewServer(rdb *redis.Client) *Server {
	server := &Server{
		Port: "8080",
		rdb:  rdb,
	}
	return server
}

func (s *Server) StartServer() {
	//start server and pass params into redis
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	//8080
	r.Run(s.Port)
}
