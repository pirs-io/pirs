package register

import (
	"pirs.io/tracker/redis"
)

type TrackerInstance struct {
	redis.ReJsonDocument `json:"-"`
	OrganizationName     string   `json:"OrganizationName,omitempty"`
	Url                  string   `json:"Url,omitempty"`
	RegisteredInstances  []string `json:"RegisteredInstances,omitempty"`
}
