package register

import "pirs.io/common"

type TrackerInstance struct {
	common.ReJsonDocument `json:"-"`
	OrganizationName      string   `json:"OrganizationName,omitempty"`
	Url                   string   `json:"Url,omitempty"`
	RegisteredInstances   []string `json:"RegisteredInstances,omitempty"`
}
