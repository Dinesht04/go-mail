package services

import (
	"context"

	"github.com/dinesht04/go-micro/internal/data"
	"github.com/redis/go-redis/v9"
)

func ValidateTask(task data.Task, rdb *redis.Client, ctx context.Context) (status bool, message string, err error) {
	if task.Type == "subscribe" || task.Type == "unsubscribe" {
		exists, err := rdb.Exists(ctx, "subscriptionContentMap"+task.Payload.ContentType).Result()
		if err != nil {
			return false, "Internal Server Error", err
		}
		if exists == 1 {
			return true, "Validation Sucess", nil
		} else {
			return false, "Content Type doesn't exist", nil
		}

	}
	return true, "Validation Sucess", nil
}
