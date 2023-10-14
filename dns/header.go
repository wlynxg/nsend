package dns

import (
	"encoding/binary"
	"errors"
	"math/rand"
)

type Header struct {
	TransactionID uint16
	Flags         uint16
	Questions     uint16
	AnswerRRs     uint16
	AuthorityRRs  uint16
	AdditionalRRs uint16
}

func DefaultHeader() Header {
	return Header{
		TransactionID: uint16(rand.Intn(1 << 16)),
		Flags:         0x0100,
		Questions:     1,
		AnswerRRs:     0,
		AuthorityRRs:  0,
		AdditionalRRs: 0,
	}
}

func MarshalHeader(h *Header) []byte {
	if h == nil {
		return nil
	}

	buff := make([]byte, 12)
	binary.BigEndian.PutUint16(buff, h.TransactionID)
	binary.BigEndian.PutUint16(buff[2:], h.Flags)
	binary.BigEndian.PutUint16(buff[4:], h.Questions)
	binary.BigEndian.PutUint16(buff[6:], h.AnswerRRs)
	binary.BigEndian.PutUint16(buff[8:], h.AuthorityRRs)
	binary.BigEndian.PutUint16(buff[10:], h.AdditionalRRs)
	return buff
}

func UnmarshalHeader(buff []byte) (*Header, error) {
	if len(buff) < 12 {
		return nil, errors.New("this is not a complete DNS packet header")
	}

	header := &Header{}
	header.TransactionID = binary.BigEndian.Uint16(buff[0:2])
	header.Flags = binary.BigEndian.Uint16(buff[2:4])
	header.Questions = binary.BigEndian.Uint16(buff[4:6])
	header.AnswerRRs = binary.BigEndian.Uint16(buff[6:8])
	header.AuthorityRRs = binary.BigEndian.Uint16(buff[8:10])
	header.AdditionalRRs = binary.BigEndian.Uint16(buff[10:12])
	return header, nil
}
