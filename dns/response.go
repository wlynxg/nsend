package dns

type Response struct {
	Header
	Queries []Query
	Answers []Answer
}

func UnmarshalResponse(buff []byte) (*Response, error) {
	return nil, nil
}
