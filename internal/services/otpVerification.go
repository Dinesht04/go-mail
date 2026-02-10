package services

import (
	"github.com/dinesht04/go-micro/internal/data"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func VerifyOtp(data data.VerifyOtpParams, rdb *redis.Client, ctx *gin.Context) (bool, error) {

	res, err := rdb.HGet(ctx, "otp_hashmap", data.UserEmail).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		} else {
			return false, err
		}
	}

	if res == data.Otp {
		return true, nil
	} else {
		return false, nil
	}

}
