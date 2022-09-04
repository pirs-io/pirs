package service

import (
	"pirs.io/pirs/common"
	"pirs.io/pirs/common/trackerProto"
	"pirs.io/pirs/tracker/domain/location"
)

var (
	lsLog = common.GetLoggerFor("location_service")
)

type LocationService struct {
	LocationRepository location.Repository
}

func (locService *LocationService) RegisterPackage(info *trackerProto.PackageInfo) {
	savePackage, err := locService.LocationRepository.SavePackage(info)
	if err != nil {
		lsLog.Fatal().Msgf("Error saving package! %s", err)
	}
	lsLog.Info().Msgf("Saved package %s", savePackage.String())
}
