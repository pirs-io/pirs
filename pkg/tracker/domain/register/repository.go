package register

type Repository interface {
	GetAllRegisteredInstances() ([]TrackerInstance, error)
	SaveTrackerNewInstanceData(peer TrackerInstance) error
	SaveNetworkRegisteredPeerData(peer TrackerInstance) error
}
