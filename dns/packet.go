package dns

const BufferSize = 512

type QueryType int

const (
	UnknownQueryType QueryType = iota
	A
	AAAA
)

func (t QueryType) String() string {
	switch t {
	case UnknownQueryType:
		return "UnknownQueryType"
	case A:
		return "A"
	case AAAA:
		return "AAAA"
	default:
		return ""
	}
}

type ClassType int

const (
	UnknownClass ClassType = iota
	IN
)

func (t ClassType) String() string {
	switch t {
	case UnknownClass:
		return "UnknownClass"
	case IN:
		return "IN"
	default:
		return ""
	}
}
