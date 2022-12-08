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
		return "PETRIFLOW"
	case BPMN:
		return "BPMN"
	}
	return "UNKNOWN"
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
