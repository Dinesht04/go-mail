package cron

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dinesht04/go-micro/internal/data"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
)

type userRecord struct {
	ClientEmail string
	cronId      cron.EntryID
}

type CronJobStation struct {
	context context.Context
	rdb     *redis.Client
	cron    *cron.Cron
	Jobs    map[string]userRecord
}

func CreateNewCronJobStation(ctx context.Context, rdb *redis.Client) *CronJobStation {
	c := cron.New()
	c.Start()
	return &CronJobStation{
		cron:    c,
		Jobs:    make(map[string]userRecord),
		context: ctx,
		rdb:     rdb,
	}
}

func (c *CronJobStation) Subscribe(userEmailId string, frequency string, contentType string) error {

	cronId, err := RegisterCronSendingEmailJob(c, userEmailId, frequency, contentType)
	if err != nil {
		return err
	}

	record := userRecord{
		cronId: cronId,
	}

	c.Jobs[userEmailId+contentType] = record
	fmt.Println("cron job added successfully")
	return nil
}

func (c *CronJobStation) Unsubscribe(userEmailId string, contentType string) error {
	Record, ok := c.Jobs[userEmailId+contentType]
	if !ok {
		return fmt.Errorf("Record doesnt exist how to unsubscruibe?")
	}
	c.cron.Remove(Record.cronId)
	delete(c.Jobs, userEmailId)
	fmt.Println("cron job removed successfully")
	return nil
}

func RegisterCronSendingEmailJob(c *CronJobStation, userEmailId string, frequency string, contentType string) (cron.EntryID, error) {
	fmt.Println("Registering for the job")

	cronId, err := c.cron.AddFunc("*/1 * * * *", func() {
		//this stuff goes to logs
		fmt.Println("Sending a mail")
		content, err := c.rdb.HGetAll(c.context, "subscriptionContentMap"+contentType).Result()
		if err != nil {
			if err == redis.Nil {
				fmt.Println("This type of content doesnt exist")
				return
			} else {
				fmt.Println("Err accessing content type")
				fmt.Println(err)
				return
			}
		}

		messageTask := data.Task{
			Id:   uuid.NewString(),
			Task: "Automated Email",
			Type: "message",
			Payload: data.Payload{
				UserID:  userEmailId,
				Subject: content["subject"],
				Content: content["content"],
			},
			Retries: 3,
		}

		encodedTask, err := json.Marshal(&messageTask)
		if err != nil {
			fmt.Println("Error decoding")
			fmt.Println(err)
			return
		}

		err = c.rdb.RPush(c.context, "taskQueue", encodedTask).Err()
		if err != nil {
			fmt.Println("Error Pushing to task Queue")
			fmt.Println(err)
			return
		} else {
			fmt.Println("Cron job added automatic email to taskqueue successully")
		}

	})

	return cronId, err
}
