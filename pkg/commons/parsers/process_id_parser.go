package parsers

import (
	"errors"
	"github.com/rs/zerolog/log"
	"path/filepath"
	"strconv"
	"strings"
)

type ProcessId struct {
	Organization string
	Tenant       string
	Project      string
	Process      string
	Version      int64
}

var (
	ErrProcessIdBadFormat = errors.New("bad format for processId. Format is: organization.tenant.project.process:version")
)

func ParseProcessId(processId string) (*ProcessId, error) {
	splitted := make([]string, 0)
	splittedWithoutVersionParsed := strings.Split(processId, ".")
	if len(splittedWithoutVersionParsed) != 4 {
		return nil, ErrProcessIdBadFormat
	}
	splitted = append(splitted, splittedWithoutVersionParsed[0:3]...)
	versionSplitted := strings.Split(splittedWithoutVersionParsed[3], ":")
	if len(versionSplitted) != 2 {
		return nil, ErrProcessIdBadFormat
	}
	splitted = append(splitted, versionSplitted...)
	version, err := strconv.Atoi(splitted[4])
	if err != nil {
		log.Err(err)
		return nil, err
	}
	return &ProcessId{
		Organization: splitted[0],
		Tenant:       splitted[1],
		Project:      splitted[2],
		Process:      splitted[3],
		Version:      int64(version),
	}, nil
}

func (r *ProcessId) ProcessWithinProject() *string {
	var res = filepath.Join(r.Tenant, r.Project, r.Process)
	return &res
}

func (r *ProcessId) FullProcessIdWithVersionTag() string {
	return strings.Join([]string{r.Organization, r.Tenant, r.Project, r.Process}, ".") + ":" + strconv.Itoa(int(r.Version))
}

func (r *ProcessId) FullProcessIdWithoutVersionTag() string {
	return strings.Join([]string{r.Organization, r.Tenant, r.Project, r.Process}, ".")
}
