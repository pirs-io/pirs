package db

import (
	"context"
	"github.com/go-redis/redis/v9"
	"pirs.io/pirs/common/trackerProto"
)

type LocationRepo struct {
	Context *context.Context
	Client  *redis.Client
}

func (r *LocationRepo) SavePackage(info *trackerProto.PackageInfo) (trackerProto.PackageInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (r *LocationRepo) GetPackageById(info *trackerProto.PackageInfo) (trackerProto.PackageInfo, error) {
	//TODO implement me
	panic("implement me")
}
