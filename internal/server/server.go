package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dinesht04/go-micro/internal/data"
	"github.com/dinesht04/go-micro/internal/email"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	Port string
	rdb  *redis.Client
}

func NewServer(rdb *redis.Client) *Server {
	server := &Server{
		Port: ":8080",
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

	r.POST("/task", func(ctx *gin.Context) {
		var task data.Task
		err := ctx.ShouldBind(&task)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"ERROR": "INVALID FORMAT",
			})
			log.Fatal(err)
		}
		fmt.Println(task)

		encodedTask, err := json.Marshal(&task)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"ERROR": "Error while marhsalling task",
			})
			return
		}

		err = s.rdb.RPush(ctx, "taskQueue", encodedTask).Err()
		if err != nil {
			log.Fatal(err)
		}

		//mantain a map in memory

		//how to implement the retries mechanic?
		//how will the queue insertion work? - draw on excalidraw

	})

	r.GET("/task", func(ctx *gin.Context) {

		//log tasks here?

		val := s.rdb.RPop(ctx, "taskQueue")

		if err := val.Err(); err == redis.Nil {
			ctx.JSON(http.StatusOK, gin.H{
				"Status": "Queue is empty",
			})
			return
		}

		var Task data.Task

		encodedTask, err := val.Bytes()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"ERROR": "Error while decoding redis string",
			})
			fmt.Println(val.Err())
			return
		}

		err = json.Unmarshal(encodedTask, &Task)

		fmt.Println("redis string: ", val.String())

		ctx.JSON(http.StatusOK, gin.H{
			"task": Task,
		})

	})

	r.GET("/verify", func(ctx *gin.Context) {
		var req data.VerifyOtpParams
		err := ctx.ShouldBind(&req)
		if err != nil {
			log.Fatal(err)
		}

		verified := email.VerifyOtp(req, s.rdb, ctx)

		ctx.JSON(http.StatusOK, gin.H{
			"type":     "otp verification",
			"verified": verified,
		})

	})

	//how to manage logs?

	r.GET("/updateSubscriptionContent", func(ctx *gin.Context) {
		// err := rdb.HSet(ctx, "subscriptionContentMap"+task.Payload.ContentType, fields).Err()

	})

	r.GET("/metrics", func(ctx *gin.Context) {
		//us redis to store and access total jobs, successful jobs, etv
	})

	r.GET("/logs", func(ctx *gin.Context) {
		//implement all that elastisearch shi?
	})

	//8080
	r.Run(s.Port)
}
