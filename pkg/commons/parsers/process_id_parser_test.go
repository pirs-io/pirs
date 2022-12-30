package parsers

import (
	"errors"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"pirs.io/commons"
	"strconv"
	"testing"
)

func TestParseProcessId(t *testing.T) {
	tests := []struct {
		name  string
		in    string
		want  ProcessId
		error commons.ErrorResponse
	}{
		{
			name: "valid processId",
			in:   "organization.tenant.project.process:1",
			want: ProcessId{
				Organization: "organization",
				Tenant:       "tenant",
				Project:      "project",
				Process:      "process",
				Version:      1,
			},
		},
		{
			name:  "missing version",
			in:    "organization.tenant.project.process",
			error: ErrProcessIdBadFormat,
		},
		{
			name:  "incorrect version",
			in:    "organization.tenant.project.process:one",
			error: &strconv.NumError{Func: "Atoi", Num: "one", Err: errors.New("invalid syntax")},
		},
		{
			name:  "not complete uri",
			in:    "organization.project.process:1",
			error: ErrProcessIdBadFormat,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			processId, err := ParseProcessId(test.in)
			if test.error != nil {
				if err.Error() != test.error.Error() {
					t.Fatalf("want: %s, got: %s", test.error.Error(), err.Error())
				}
				return
			}
			if &test.want != nil {
				if !cmp.Equal(processId, &test.want) {
					t.Fatalf("want: %s, got: %s", fmt.Sprint(test.want), fmt.Sprint(processId))
				}
				return
			}
		})
	}
}
