package db

import (
	"context"
	"pirs.io/common"
	"pirs.io/tracker/domain/register"
)

const (
	keyRegisteredInstances = "instance:"
)

var (
	rrLog = common.GetLoggerFor("register_repo_logger")
)

type RegisterRepo struct {
	Context *context.Context
	Client  *common.CustomRedisClient
}

func (r *RegisterRepo) GetAllRegisteredInstances() ([]register.TrackerInstance, error) {
	panic("Not implemented!")
}

func (r *RegisterRepo) SaveTrackerNewInstanceData(peer register.TrackerInstance) error {
	panic("Not implemented!")

}

func (r *RegisterRepo) SaveNetworkRegisteredPeerData(peer register.TrackerInstance) error {
	panic("Not implemented!")
}

func makeKey(peer register.TrackerInstance) string {
	return keyRegisteredInstances + peer.OrganizationName
}
