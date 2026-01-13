package main

import (
	"context"
	"fmt"
	"log"
	"os"

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

	smtpPass := os.Getenv("smtp_pass")
	fmt.Println(smtpPass)

	ctx := context.Background()

	rdb := data.NewRedisClient(ctx)
	server := server.NewServer(rdb)

	Workstation := worker.NewWorkStation(rdb, 3)
	Workstation.StartWorkers(ctx)

	server.StartServer()

}
