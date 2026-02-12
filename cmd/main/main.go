package main

import (
	"context"
	"fmt"

	"github.com/dinesht04/go-micro/internal/cron"
	"github.com/dinesht04/go-micro/internal/data"
	"github.com/dinesht04/go-micro/internal/log"
	"github.com/dinesht04/go-micro/internal/server"
	"github.com/dinesht04/go-micro/internal/worker"
	"github.com/joho/godotenv"
)

func main() {

	logger, file, err := log.CreateLogger()
	if err != nil {
		panic(fmt.Errorf("Error craeting Logger"))
	}
	defer file.Close()

	err = godotenv.Load()
	if err != nil {
		logger.Error("Error Loading .env file", "error", err)
		panic(fmt.Errorf("Error Loading .env file"))
	}

	ctx := context.Background()

	rdb, err := data.NewRedisClient(ctx, logger)
	if err != nil {
		logger.Error("Error Initiating redis client", "error", err)
	}

	server := server.NewServer(rdb, logger)
	CronJobStation := cron.CreateNewCronJobStation(ctx, rdb, logger)

	Workstation := worker.NewWorkStation(rdb, 3, CronJobStation)
	Workstation.StartWorkers(ctx, logger)

	server.StartServer()

}
