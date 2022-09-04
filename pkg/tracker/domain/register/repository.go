package register

type Repository interface {
	GetAllRegisteredInstances() ([]TrackerInstance, error)
	RegisterInstance(peer TrackerInstance) error
}
