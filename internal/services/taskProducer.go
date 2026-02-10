package services

import (
	"encoding/json"
	"fmt"

	"github.com/dinesht04/go-micro/internal/data"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func ProduceTask(task data.Task, rdb *redis.Client, ctx *gin.Context) (bool, string, error) {
	task.Id = uuid.NewString()

	encodedTask, err := json.Marshal(&task)
	if err != nil {
		return false, "Error while marshalling task", err
	}

	err = rdb.RPush(ctx, "taskQueue", encodedTask).Err()
	if err != nil {
		return false, "Redis Error while pushing task to the Queue", err

	} else {
		return true, fmt.Sprintf("Pushed task to the Queue Successfully, \n Task ID: %v, Task Name: %v", task.Id, task.Task), nil
	}
}
