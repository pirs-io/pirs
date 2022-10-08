package service

import (
	"pirs.io/common"
	"pirs.io/tracker/domain/register"
)

var (
	irLog = common.GetLoggerFor("location_service")
)

type InstanceRegistrationService struct {
	RegisterRepo register.Repository
}

func (r *InstanceRegistrationService) RegisterInstance(orgId string, orgUrl string) error {
	panic("Implement me!")
}

func (r *InstanceRegistrationService) SaveRegisteredPeerData() error {
	panic("Implement me!")
}

func (r *InstanceRegistrationService) GetAllRegisteredInstances() ([]register.TrackerInstance, error) {
	panic("Implement me!")
}
