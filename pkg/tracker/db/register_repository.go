package db

import (
	"context"
	"github.com/samber/lo"
	"pirs.io/pirs/common"
	"pirs.io/pirs/tracker/domain/register"
)

const (
	keyRegisteredInstances = "INSTANCES"
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

func (r *RegisterRepo) RegisterInstance(peer register.TrackerInstance) error {
	// check all existing organizations
	existingOrganizations, err := r.Client.ObjKeys(keyRegisteredInstances, "$")
	if err != nil {
		rrLog.Warn().Msg(err.Error())
	}
	if lo.Contains[string](existingOrganizations, peer.OrganizationName) {
		rrLog.Warn().Msgf("Organization %s already registered!", peer.OrganizationName)
		return nil
	}
	// if structure for instances is not created - initialize it
	if existingOrganizations == nil || len(existingOrganizations) == 0 {
		_, err := r.Client.JsonSetString(keyRegisteredInstances, "$", "{}")
		if err != nil {
			return err
		}
	}
	// set entry for organization registering to this instance
	_, err = r.Client.JsonSet(keyRegisteredInstances, "$."+peer.OrganizationName, peer)
	if err != nil {
		rrLog.Error().Msg(err.Error())
	}
	return nil
}
