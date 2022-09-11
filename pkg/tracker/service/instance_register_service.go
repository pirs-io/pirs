package service

import (
	"github.com/rs/zerolog/log"
	"pirs.io/pirs/common"
	"pirs.io/pirs/common/trackerProto"
	"pirs.io/pirs/tracker/domain/register"
)

var (
	irLog = common.GetLoggerFor("location_service")
)

type InstanceRegistrationService struct {
	RegisterRepo register.Repository
}

func (r *InstanceRegistrationService) RegisterInstance(info *trackerProto.TrackerInfo) (*trackerProto.InstanceRegisterResponse, error) {
	irLog.Info().Msgf("Registering %s", info.OrganizationId)
	err := r.RegisterRepo.RegisterInstance(register.TrackerInstance{
		OrganizationName: info.OrganizationId,
		Url:              "localhost:8081",
	})
	if err != nil {
		log.Err(err)
		return &trackerProto.InstanceRegisterResponse{
			Status: trackerProto.InstanceRegisterStatus_FAILED,
			Error:  "",
		}, err
	} else {
		return &trackerProto.InstanceRegisterResponse{
			Status: trackerProto.InstanceRegisterStatus_SUCCESS,
			Error:  "",
		}, nil
	}
}

func (r *InstanceRegistrationService) GetAllRegisteredInstances() ([]register.TrackerInstance, error) {
	return r.RegisterRepo.GetAllRegisteredInstances()
}
