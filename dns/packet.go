package dns

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

type Header struct {
	TransactionID int16
	Flags         int16
	Questions     int16
	AnswerRRs     int16
	AuthorityRRs  int16
	AdditionalRRs int16
}

type Request struct {
	Header
	Queries []Query
}

type Query struct {
	Name  string
	Type  QueryType
	Class ClassType
}

type Response struct {
	Header
	Queries []Query
	Answers []Answer
}

type Answer interface{}
