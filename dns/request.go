package dns

type Request struct {
	Header
	Queries []Query
}

func MarshalRequest(r *Request) []byte {
	if r == nil {
		return nil
	}

	var (
		buff   = make([]byte, BufferSize)
		offset = 0
	)

	header := MarshalHeader(&r.Header)
	copy(buff[:], header)
	offset += len(header)

	for _, query := range r.Queries {
		data := MarshalQuery(&query)
		copy(buff[offset:], data)
		offset += len(data)
	}
	return buff[:offset]
}

func NewRequest(domain string, queryType QueryType) *Request {
	req := &Request{
		Header: DefaultHeader(),
		Queries: []Query{{
			Name:  domain,
			Type:  queryType,
			Class: IN,
		}},
	}
	return req
}
