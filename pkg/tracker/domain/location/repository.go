package location

import "pirs.io/pirs/common/trackerProto"

type Repository interface {
	SavePackage(info *trackerProto.PackageInfo) (trackerProto.PackageInfo, error)
	GetPackageById(info *trackerProto.PackageInfo) (trackerProto.PackageInfo, error)
}
