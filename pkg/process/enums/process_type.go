package enums

type ProcessType int

const (
	Petriflow ProcessType = iota
	BPMN
)

func (pt ProcessType) String() string {
	switch pt {
	case Petriflow:
		return "petriflow"
	case BPMN:
		return "bpmn"
	}
	return "unknown"
}
