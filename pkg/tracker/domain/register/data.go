package register

type TrackerInstance struct {
	OrganizationName string `json:"OrganizationName,omitempty"`
	Url              string `json:"Url,omitempty"`
}

func (r TrackerInstance) isDocument() {}
