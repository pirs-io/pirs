package common

import "github.com/go-redis/redis/v9"

func NewRedisClient(addr string, port string, pwd string, db int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr + ":" + port,
		Password: pwd,
		DB:       db,
	})
	return rdb
}
