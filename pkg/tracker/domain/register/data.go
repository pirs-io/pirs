package register

import "pirs.io/pirs/common"

type TrackerInstance struct {
	common.ReJsonDocument `json:"-"`
	OrganizationName      string `json:"OrganizationName,omitempty"`
	Url                   string `json:"Url,omitempty"`
}
