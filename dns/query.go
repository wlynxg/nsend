package dns

import (
	"bytes"
	"encoding/binary"
	"errors"
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

func UnmarshalQuery(raw []byte) (*Query, error) {
	if len(raw) < 6 {
		return nil, errors.New("this is not a valid queries slice")
	}

	var (
		offset int
		domain bytes.Buffer
		query  *Query
	)

	for raw[offset] != 0 {
		length := int(raw[offset])
		offset++
		domain.Write(raw[offset : offset+length])
		domain.WriteByte('.')
		offset += length
	}
	offset++
	query.Name = domain.String()[:domain.Len()-1]
	query.Type = QueryType(binary.BigEndian.Uint16(raw[offset : offset+2]))
	offset += 2
	query.Class = ClassType(binary.BigEndian.Uint16(raw[offset : offset+2]))
	offset += 2
	return query, nil
}
