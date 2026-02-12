package worker

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"time"

	"github.com/dinesht04/go-micro/internal/cron"
	"github.com/dinesht04/go-micro/internal/data"
	"github.com/dinesht04/go-micro/internal/services"
	"github.com/redis/go-redis/v9"
)

type WorkStation struct {
	Rdb         *redis.Client
	Workers     int
	CronStation *cron.CronJobStation
}

func NewWorkStation(rdb *redis.Client, num int, cron *cron.CronJobStation) *WorkStation {
	return &WorkStation{
		Rdb:         rdb,
		Workers:     num,
		CronStation: cron,
	}
}

func (w *WorkStation) StartWorkers(ctx context.Context, logger *slog.Logger) {

	for range w.Workers {
		go Worker(w.Rdb, ctx, w.CronStation, logger)
	}
}

func Worker(rdb *redis.Client, ctx context.Context, cron *cron.CronJobStation, logger *slog.Logger) {

	for {

		results, err := rdb.BLPop(ctx, time.Minute, "taskQueue").Result()
		if err != nil {
			if err == redis.Nil {
				continue
			} else {
				logger.Error("Error Popping Task from Queue",
					"error", err,
				)

			}
		}

		result := results[1]

		var task data.Task

		err = json.Unmarshal([]byte(result), &task)
		if err != nil {
			logger.Error("Error Unmarshalling task",
				"error", err,
				"taskId", task.Id,
				"taskName", task.Task)
			continue
		}

		task.Retries = task.Retries - 1

		logger.Info("Task Popped from queue",
			"taskId", task.Id,
			"taskName", task.Task)

		taskType := task.Type

		var status bool
		var logs string

		err = rdb.Incr(ctx, "totalTasksExecuted").Err()
		if err != nil {
			logger.Error("Error incrementing total tasks", "error", err)

		}

		switch taskType {
		case "generateOtp":
			status, logs, err = services.GenerateOtp(task, rdb, ctx)
		case "message":
			//this can stay here
			status, logs, err = services.Sendmessage(task, rdb)
		case "subscribe":
			//This can stay here
			status, logs, err = services.Subscribe(task, rdb, ctx, cron)
		case "unsubscribe":
			//should this stay here?
			status, logs, err = services.Unsubscribe(task, rdb, cron)
		default:
			logger.Warn("Invalid Task Type", "unknown_task_type", task.Type)
		}

		logger.Info("Task processed",
			"log", logs,
			"taskId", task.Id,
			"taskName", task.Task,
			"taskType", taskType)

		// error or status pe check?
		if !status {

			err = rdb.Incr(ctx, "totalTasksFailed").Err()
			if err != nil {
				logger.Error("Error incrementing failed tasks", "error", err)

			}

			if err != nil {
				logger.Error("Task Failed Due to error",
					"error", err,
					"taskId", task.Id,
					"taskName", task.Task,
					"taskType", taskType,
					"Retries left", task.Retries,
				)

			} else {
				logger.Warn("Task Failed",
					"latest_logs", logs,
					"taskId", task.Id,
					"taskName", task.Task,
					"taskType", taskType,
					"Retries left", task.Retries,
				)

			}

			if task.Retries <= 0 {
				logger.Warn("Retries Finished",
					"latest_logs", logs,
					"taskId", task.Id,
					"taskName", task.Task,
					"taskType", taskType,
					"Retries left", task.Retries)
				continue
			}

			logger.Info("Adding task back to Queue...",
				"latest_logs", logs,
				"taskId", task.Id,
				"taskName", task.Task,
				"taskType", taskType,
				"Retries left", task.Retries)

			encodedTask, err := json.Marshal(&task)
			if err != nil {
				log.Fatal(err)
			}

			err = rdb.RPush(ctx, "taskQueue", encodedTask).Err()
			if err != nil {
				log.Fatal(err)
			}

		} else {
			logger.Info("Performed Task Successfully!!!",
				"latest_logs", logs,
				"taskId", task.Id,
				"taskName", task.Task,
				"taskType", taskType)

			err = rdb.Incr(ctx, "totalTasksSuccessful").Err()
			if err != nil {
				logger.Error("Error incrementing successful tasks", "error", err)

			}
		}

	}

}
