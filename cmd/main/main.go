package main

import (
	"context"
	"log"

	"github.com/dinesht04/go-micro/internal/cron"
	"github.com/dinesht04/go-micro/internal/data"
	"github.com/dinesht04/go-micro/internal/server"
	"github.com/dinesht04/go-micro/internal/worker"
	"github.com/joho/godotenv"
)

func main() {

	//connect to redis
	//start the server

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error Loading .env file")
	}

	ctx := context.Background()

	rdb := data.NewRedisClient(ctx)
	server := server.NewServer(rdb)
	CronJobStation := cron.CreateNewCronJobStation(ctx, rdb)

	Workstation := worker.NewWorkStation(rdb, 3, CronJobStation)
	Workstation.StartWorkers(ctx)

	server.StartServer()

}
