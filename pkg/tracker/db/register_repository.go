package db

import (
	"context"
	"github.com/go-redis/redis/v9"
	"pirs.io/pirs/common"
	"pirs.io/pirs/tracker/domain/register"
	"time"
)

var (
	rrLog = common.GetLoggerFor("register_repo_logger")
)

type RegisterRepo struct {
	Context *context.Context
	Client  *redis.Client
}

func (r *RegisterRepo) GetAllRegisteredInstances() ([]register.TrackerInstance, error) {
	return []register.TrackerInstance{}, nil
}

func (r *RegisterRepo) RegisterInstance(peer register.TrackerInstance) error {
	rrLog.Debug().Msgf("Registering instance: %s", peer.OrganizationName)
	r.Client.Set(*r.Context, register.KEY_REGISTERED_INSTANCES, "v", time.Duration(0))
	return nil
}
