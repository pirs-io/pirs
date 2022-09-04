package common

import (
	"context"
	"github.com/go-redis/redis/v9"
)

func NewRedisClient(ctx context.Context, addr string, port string, pwd string, db int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr + ":" + port,
		Password: pwd,
		DB:       db,
	})
	pong := rdb.Ping(ctx)
	if pong.Err() != nil {
		panic(pong.Err())
	}
	return rdb
}
