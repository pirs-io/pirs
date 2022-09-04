package db

import (
	"github.com/go-redis/redis/v9"
	"pirs.io/pirs/common"
	"pirs.io/pirs/tracker/domain/register"
)

var (
	rrLog = common.GetLoggerFor("register_repo_logger")
)

type RegisterRepo struct {
	Client *redis.Client
}

func (r *RegisterRepo) GetAllRegisteredInstances() ([]register.TrackerInstance, error) {
	return []register.TrackerInstance{}, nil
}

func (r *RegisterRepo) RegisterInstance(peer register.TrackerInstance) error {
	rrLog.Debug().Msgf("Registering instance: %s", peer.OrganizationName)
	return nil
}
