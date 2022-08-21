package db

import (
	"github.com/go-redis/redis/v9"
	"pirs.io/pirs/common/trackerProto"
)

type RedisRepo struct {
	Client *redis.Client
}

func (r *RedisRepo) SavePackage(info *trackerProto.PackageInfo) (trackerProto.PackageInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RedisRepo) GetPackageById(info *trackerProto.PackageInfo) (trackerProto.PackageInfo, error) {
	//TODO implement me
	panic("implement me")
}
