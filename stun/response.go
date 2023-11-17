package stun

import (
	"encoding/binary"
	"net/netip"

	"github.com/pkg/errors"
)

type Response struct {
	MessageType   int
	MessageLength int
	MagicCookie   uint32
	TransactionID TxID
	Attributes    map[AttributeType]Attribute
}

func UnmarshalResponse(buff []byte, resp *Response) (int, error) {
	if len(buff) < 20 {
		return 0, errors.New("invalid stun response packet")
	}

	var offset = 0

	if resp.Attributes == nil {
		resp.Attributes = make(map[AttributeType]Attribute)
	}

	// set the MessageType
	if binary.BigEndian.Uint16(buff[offset:offset+2]) == BindingResponse {
		resp.MessageType = BindingResponse
	}
	offset += 2

	// set the MessageLength
	resp.MessageLength = int(binary.BigEndian.Uint16(buff[offset : offset+2]))
	offset += 2

	// set the MagicCookie
	resp.MagicCookie = binary.BigEndian.Uint32(buff[offset : offset+4])
	offset += 4

	// set the TransactionID
	resp.TransactionID = make(TxID, 12)
	copy(resp.TransactionID, buff[offset:offset+12])
	offset += 12

	old := offset
	for offset-old < resp.MessageLength {
		attribute := Attribute{}

		// set AttributeType
		attribute.Type = AttributeType(binary.BigEndian.Uint16(buff[offset : offset+2]))
		offset += 2

		// set AttributeLength
		attribute.Length = int(binary.BigEndian.Uint16(buff[offset : offset+2]))
		offset += 2

		// don't parse comprehension-option
		if attribute.Type.IsComprehensionOptional() {
			offset += attribute.Length
			offset += attribute.Length % Padding
		}

		// set AttributeReserved
		attribute.Reserved = int(buff[offset])
		offset += 1

		// set ProtocolFamily
		attribute.ProtocolFamily = ProtocolFamily(buff[offset])
		offset += 1

		// set Port
		attribute.Port = int(binary.BigEndian.Uint16(buff[offset : offset+2]))
		offset += 2

		// set IP
		attribute.IP = netip.AddrFrom4([4]byte{buff[offset], buff[offset+1], buff[offset+2], buff[offset+3]})
		offset += 4

		resp.Attributes[attribute.Type] = attribute
	}

	return offset, nil
}
