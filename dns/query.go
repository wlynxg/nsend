package dns

import (
	"encoding/binary"
	"strings"
)

type Query struct {
	Name  string
	Type  QueryType
	Class ClassType
}

/*
raw: 06 67 6f 6f 67 6c 65 | 03 63 6f 6d | 00 | 00 01 | 00 01                                     ..
->:  6  g  o  o  g  l  e  | 3  c  o  m  | 0    A       IN
	 ^                      ^  		      ^
len("google")           len("com")       EOF
*/

func MarshalQuery(q *Query) []byte {
	var (
		buff   = make([]byte, BufferSize)
		offset = 0
	)

	domains := strings.Split(q.Name, ".")
	for _, domain := range domains {
		buff[offset] = byte(len(domain))
		offset += 1
		copy(buff[offset:], domain)
		offset += len(domain)
	}
	buff[offset] = 0
	offset += 1

	binary.BigEndian.PutUint16(buff[offset:], uint16(q.Type))
	offset += 2
	binary.BigEndian.PutUint16(buff[offset:], uint16(q.Class))
	offset += 2

	return buff[:offset]
}
