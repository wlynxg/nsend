package dns

type Response struct {
	Header
	Queries []Query
	Answers []Answer
}

type Answer interface{}

func NewResponse(buff []byte) *Response {

	return nil
}
