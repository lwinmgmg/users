package services

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	UserRedis *redis.Client = nil
)

func GetKey(key string) (string, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		1000*time.Millisecond,
	)
	defer cancel()
	res, err := UserRedis.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return res, nil
}

func SetKey(key, val string, duration time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		1000*time.Millisecond,
	)
	defer cancel()
	res, err := UserRedis.Set(ctx, key, val, duration).Result()
	if err != nil {
		return "", err
	}
	return res, nil
}
