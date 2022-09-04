package db

import (
	"github.com/go-redis/redis/v9"
	"pirs.io/pirs/common/trackerProto"
)

type LocationRepo struct {
	Client *redis.Client
}

func (r *LocationRepo) SavePackage(info *trackerProto.PackageInfo) (trackerProto.PackageInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (r *LocationRepo) GetPackageById(info *trackerProto.PackageInfo) (trackerProto.PackageInfo, error) {
	//TODO implement me
	panic("implement me")
}
