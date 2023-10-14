package dns

const BufferSize = 512

type QueryType int

const (
	UnknownQueryType QueryType = iota
	A
	AAAA
)

type ClassType int

const (
	UnknownClass ClassType = iota
	IN
)
