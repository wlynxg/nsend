package dns

import (
	"math/rand"

	"github.com/wlynxg/nsend/pkg/xmap"
)

var cache xmap.Map[int16, struct{}]

func NewRequest(domain string, queryType QueryType) *Request {
	tid := int16(rand.Intn(1 << 16))
	for cache.CompareAndDelete(tid, struct{}{}) {
		tid = int16(rand.Intn(1 << 16))
	}

	req := &Request{
		Header: Header{
			TransactionID: tid,
			Flags:         0x0100,
			Questions:     1,
			AnswerRRs:     0,
			AuthorityRRs:  0,
			AdditionalRRs: 0,
		},
		Queries: []Query{{
			Name:  domain,
			Type:  queryType,
			Class: IN,
		}},
	}
	return req
}
