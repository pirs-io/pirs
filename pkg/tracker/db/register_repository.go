package db

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v9"
	"pirs.io/pirs/common"
	"pirs.io/pirs/tracker/domain/register"
	"strings"
)

var (
	rrLog = common.GetLoggerFor("register_repo_logger")
)

type RegisterRepo struct {
	Context *context.Context
	Client  *redis.Client
}

func (r *RegisterRepo) GetAllRegisteredInstances() ([]register.TrackerInstance, error) {
	keys, err := r.Client.Do(*r.Context, "JSON.OBJKEYS", "INSTANCES", "$").Result()
	if err != nil {
		return nil, err
	}
	parsedKeys := strings.Replace(strings.TrimSuffix(strings.TrimPrefix(fmt.Sprint(keys), "["), "]"), " ", ",", 1)
	rrLog.Info().Msgf("Found registered organizations: %s", parsedKeys)
	//existingInstances, err := r.Client.Do(*r.Context, "JSON.GET", register.KEY_REGISTERED_INSTANCES, "$"+"'"+parsedKeys+"'").Result()
	//instances := &[]register.TrackerInstance{}
	//err = json.Unmarshal(existingInstances.([]byte), instances)
	//if err != nil {
	//	return nil, err
	//}
	return nil, nil
}

func (r *RegisterRepo) RegisterInstance(peer register.TrackerInstance) error {
	rrLog.Debug().Msgf("Registering instance: %s", peer.OrganizationName)
	serialized, err := json.Marshal(peer)
	if err != nil {
		panic(err)
	}

	res, err := r.Client.Do(*r.Context, "JSON.OBJKEYS", register.KEY_REGISTERED_INSTANCES, "$").Result()
	if err != nil {
		rrLog.Error().Msg(err.Error())
	}
	if res == nil {
		r.Client.Do(*r.Context, "JSON.SET", register.KEY_REGISTERED_INSTANCES, "$", "{}")
	}
	res, err = r.Client.Do(*r.Context, "JSON.SET", register.KEY_REGISTERED_INSTANCES, "$."+peer.OrganizationName, string(serialized)).Result()
	if err != nil {
		rrLog.Error().Msg(err.Error())
	}
	return nil
}
