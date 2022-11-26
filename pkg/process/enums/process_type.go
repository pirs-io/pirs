package enums

type ProcessType int

const (
	Petriflow ProcessType = iota
	BPMN
	UNKNOWN
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

func (pt ProcessType) Int() int {
	switch pt {
	case Petriflow:
		return 0
	case BPMN:
		return 1
	}
	return -1
}
