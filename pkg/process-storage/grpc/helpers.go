package grpc

import (
	"errors"
	"strings"
)

type ProcessId struct {
	Organization string
	Tenant       string
	Project      string
	Process      string
}

var (
	ErrProcessIdBadFormat = errors.New("bad format for processId. Format is: organization.tenant.project.process")
)

func ParseProcessId(processId string) (*ProcessId, error) {
	splitted := strings.Split(processId, ".")
	if len(splitted) != 4 {
		return nil, ErrProcessIdBadFormat
	}
	return &ProcessId{
		Organization: splitted[0],
		Tenant:       splitted[1],
		Project:      splitted[2],
		Process:      splitted[3],
	}, nil
}

func (r *ProcessId) ProcessWithinProject() *string {
	var res = r.Project + "/" + r.Process
	return &res
}

func (r *ProcessId) FullProcessId() string {
	return strings.Join([]string{r.Organization, r.Tenant, r.Project, r.Process}, ".")
}
