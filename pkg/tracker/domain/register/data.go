package register

var (
	KEY_REGISTERED_INSTANCES = "INSTANCES"
)

type TrackerInstance struct {
	OrganizationName string
	Url              string
}
