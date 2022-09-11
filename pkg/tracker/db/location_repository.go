package db

import (
	"context"
	"pirs.io/pirs/common"
	"pirs.io/pirs/common/trackerProto"
)

type LocationRepo struct {
	Context *context.Context
	Client  *common.CustomRedisClient
}

func (r *LocationRepo) SavePackage(info *trackerProto.PackageInfo) (trackerProto.PackageInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (r *LocationRepo) GetPackageById(info *trackerProto.PackageInfo) (trackerProto.PackageInfo, error) {
	//TODO implement me
	panic("implement me")
}
