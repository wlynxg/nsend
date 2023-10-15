package dns

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net/netip"
)

type Answer struct {
	Name    string
	Type    QueryType
	Class   ClassType
	TTL     int
	Length  int
	Address netip.Addr
}

func UnmarshalAnswer(buff []byte, answer *Answer) (int, error) {
	if answer == nil {
		return -1, errors.New("answer cannot be nil")
	}

	var (
		offset int
		domain bytes.Buffer
	)

	// The dns reply enables message compression
	// https://www.rfc-editor.org/rfc/rfc1035#section-4.1.4
	if buff[0] == 0xc0 && buff[1] == 0x0c {
		answer.Name = "0xc00c"
		offset += 2
	} else {
		for buff[offset] != 0 {
			length := int(buff[offset])
			offset++
			domain.Write(buff[offset : offset+length])
			domain.WriteByte('.')
			offset += length
		}
		offset++
		answer.Name = domain.String()[:domain.Len()-1]
	}

	answer.Type = QueryType(binary.BigEndian.Uint16(buff[offset : offset+2]))
	offset += 2
	answer.Class = ClassType(binary.BigEndian.Uint16(buff[offset : offset+2]))
	offset += 2
	answer.TTL = int(binary.BigEndian.Uint32(buff[offset : offset+4]))
	offset += 4
	answer.Length = int(binary.BigEndian.Uint16(buff[offset : offset+2]))
	offset += 2

	if answer.Type == A && answer.Length == 4 {
		answer.Address = netip.AddrFrom4([4]byte{buff[offset], buff[offset+1], buff[offset+2], buff[offset+3]})
	}
	offset += answer.Length
	return offset, nil
}
