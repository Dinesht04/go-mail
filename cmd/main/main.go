package main

import (
	"context"

	"github.com/dinesht04/go-micro/internal/data"
	"github.com/dinesht04/go-micro/internal/server"
)

func main() {

	//connect to redis
	//start the server

	ctx := context.Background()

	rdb := data.NewRedisClient(ctx)
	server := server.NewServer(rdb)

	server.StartServer()

}
