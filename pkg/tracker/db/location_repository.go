package db

import (
	"github.com/go-redis/redis/v9"
	"pirs.io/pirs/common/trackerProto"
)

type RedisRepo struct {
	Client *redis.Client
}

func NewRedisClient(addr string, port string, pwd string, db int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr + ":" + port,
		Password: pwd,
		DB:       db,
	})
	return rdb
}

func (r *RedisRepo) SavePackage(info *trackerProto.PackageInfo) (trackerProto.PackageInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RedisRepo) GetPackageById(info *trackerProto.PackageInfo) (trackerProto.PackageInfo, error) {
	//TODO implement me
	panic("implement me")
}
