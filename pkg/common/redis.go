package common

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v9"
	"strings"
)

import (
	"context"
)

const (
	reJsonObjKeys = "JSON.OBJKEYS"
	reJsonObjGet  = "JSON.GET"
	reJsonObjSet  = "JSON.SET"
)

type CustomRedisClient struct {
	Client  *redis.Client
	Context context.Context
}

type ReJsonSupport interface {
	ObjKeys(key string, jsonPath string) ([]string, error)
	JsonGet(key string, jsonPath string, resList *[]ReJsonDocument) error
	JsonSet(key string, jsonPath string, data ReJsonDocument) (interface{}, error)
	JsonSetString(key string, jsonPath string, data string) (interface{}, error)
}

func NewRedisClient(ctx context.Context, addr string, port string, pwd string, db int) *CustomRedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr + ":" + port,
		Password: pwd,
		DB:       db,
	})
	pong := rdb.Ping(ctx)
	if pong.Err() != nil {
		panic(pong.Err())
	}
	return &CustomRedisClient{rdb, ctx}
}

func (c *CustomRedisClient) ObjKeys(key string, jsonPath string) ([]string, error) {
	keyRes, err := c.Client.Do(c.Context, reJsonObjKeys, key, jsonPath).Result()
	if err != nil {
		return nil, err
	}
	keys := strings.Split(strings.ReplaceAll(strings.ReplaceAll(fmt.Sprint(keyRes), "[", ""), "]", ""), " ")
	var res = make([]string, 0)
	for _, key := range keys {
		res = append(res, key)
	}
	return res, nil

}

func (c *CustomRedisClient) JsonGet(key string, jsonPath string, resList *[]ReJsonDocument) error {
	res, err := c.Client.Do(c.Context, reJsonObjGet, key, jsonPath).Result()
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(fmt.Sprint(res)), &resList)
	return err
}

func (c *CustomRedisClient) JsonSet(key string, jsonPath string, data ReJsonDocument) (interface{}, error) {
	serialized, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(c.Context, reJsonObjSet, key, jsonPath, string(serialized)).Result()
}

func (c *CustomRedisClient) JsonSetString(key string, jsonPath string, data string) (interface{}, error) {
	return c.Client.Do(c.Context, reJsonObjSet, key, jsonPath, data).Result()
}
