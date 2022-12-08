package db

import (
	"context"
	"pirs.io/commons"
	"pirs.io/tracker/domain/register"
	"pirs.io/tracker/redis"
)

const (
	keyRegisteredInstances = "instance:"
)

var (
	rrLog = commons.GetLoggerFor("register_repo_logger")
)

type RegisterRepo struct {
	Context *context.Context
	Client  *redis.CustomRedisClient
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
