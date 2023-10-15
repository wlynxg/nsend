package dns

import "errors"

type Response struct {
	Header
	Queries []Query
	Answers []Answer
}

func UnmarshalResponse(buff []byte, response *Response) (int, error) {
	if response == nil {
		return -1, errors.New("response cannot be nil")
	}

	var (
		offset = 0
	)

	length, err := UnmarshalHeader(buff[offset:], &response.Header)
	if err != nil {
		return -1, err
	}
	offset += length

	var domain string
	for i := 0; i < int(response.Questions); i++ {
		tmp := &Query{}
		length, err = UnmarshalQuery(buff[offset:], tmp)
		if err != nil {
			return -1, err
		}
		domain = tmp.Name
		offset += length
		response.Queries = append(response.Queries, *tmp)
	}

	for i := 0; i < int(response.AnswerRRs); i++ {
		tmp := &Answer{}
		length, err = UnmarshalAnswer(buff[offset:], tmp)
		if err != nil {
			return -1, err
		}
		offset += length
		if tmp.Name == "0xc00c" {
			tmp.Name = domain
		}

		response.Answers = append(response.Answers, *tmp)
	}

	return offset, nil
}
