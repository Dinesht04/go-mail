package services

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func GetMetrics(rdb *redis.Client, ctx *gin.Context) (totalTasksExecuted string, totalTasksFailed string, totalTasksSuccessful string, msg string, err error) {
	totalTasksExecuted, err = rdb.Get(ctx, "totalTasksExecuted").Result()
	if err != nil {
		if err == redis.Nil {
			totalTasksExecuted = "0"
		} else {
			return "0", "0", "0", "Error while accessing total tasks executed", err
		}
	}
	totalTasksFailed, err = rdb.Get(ctx, "totalTasksFailed").Result()
	if err != nil {
		if err == redis.Nil {
			totalTasksFailed = "0"
		} else {
			return "0", "0", "0", "Error while accessing total tasks failed", err
		}

	}
	totalTasksSuccessful, err = rdb.Get(ctx, "totalTasksSuccessful").Result()
	if err != nil {
		if err == redis.Nil {
			totalTasksSuccessful = "0"
		} else {
			return "0", "0", "0", "Error while accessing total successful tasks", err
		}
	}

	return totalTasksExecuted, totalTasksFailed, totalTasksSuccessful, "", nil

}
