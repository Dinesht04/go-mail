package services

import (
	"github.com/dinesht04/go-micro/internal/data"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func CreateContentType(data data.CreateContent, rdb *redis.Client, ctx *gin.Context) (string, error) {
	fields := []string{
		"subject", data.Subject,
		"content", data.Content,
	}

	err := rdb.HSet(ctx, "subscriptionContentMap"+data.ContentType, fields).Err()
	if err != nil {
		return "Error while updating subscription content Hashmap", err
	} else {
		return "Content Type created succesfully!", nil
	}
}

func UpdateContentType(data data.UpdateContent, rdb *redis.Client, ctx *gin.Context) (bool, string, error) {

	exists, err := rdb.Exists(ctx, "subscriptionContentMap"+data.ContentType).Result()
	if err != nil {
		return false, "Error while updating subscription content map", err
	}

	if exists == 0 {
		return false, "Content Type doesn't exist", nil
	}

	fields := []string{
		"subject", data.Subject,
		"content", data.Content,
	}

	err = rdb.HSet(ctx, "subscriptionContentMap"+data.ContentType, fields).Err()
	if err != nil {
		return false, "Error while updating subscription content map", err
	} else {
		return true, "Content map updated succesfully", nil
	}

}
