package icmp

import (
	"encoding/binary"

	"github.com/pkg/errors"
)

type PacketType byte

const (
	EchoReply   PacketType = 0
	EchoRequest PacketType = 8
)

const (
	HeaderSize = 8
)

type Packet struct {
	Type       PacketType
	Code       byte
	Checksum   uint16
	Identifier uint16
	Sequence   uint16
	Data       []byte
}

func MarshalPacket(p *Packet) []byte {
	if p == nil {
		return nil
	}

	buff := make([]byte, HeaderSize+len(p.Data))
	buff[0] = byte(p.Type)
	buff[1] = p.Code
	binary.BigEndian.PutUint16(buff[2:4], p.Identifier)
	binary.BigEndian.PutUint16(buff[6:HeaderSize], p.Sequence)
	copy(buff[HeaderSize:], p.Data)
	p.Checksum = checksum(buff)
	binary.BigEndian.PutUint16(buff[4:6], p.Checksum)
	return buff
}

func UnmarshalPacket(buff []byte, p *Packet) (int, error) {
	if p == nil {
		return -1, errors.New("response cannot be nil")
	}

	if len(buff) < HeaderSize {
		return -1, errors.New("not a correct icmp package")
	}

	p.Type = PacketType(buff[0])
	p.Code = buff[1]
	p.Checksum = binary.BigEndian.Uint16(buff[2:4])
	p.Identifier = binary.BigEndian.Uint16(buff[4:6])
	p.Sequence = binary.BigEndian.Uint16(buff[6:HeaderSize])
	if len(buff) > HeaderSize {
		data := make([]byte, len(buff)-HeaderSize)
		copy(data, buff[HeaderSize:])
		p.Data = data
	}

	return len(buff), nil
}

func NewRequest(data []byte) *Packet {
	p := &Packet{
		Type:       EchoRequest,
		Code:       0,
		Checksum:   0,
		Identifier: 0,
		Sequence:   0,
		Data:       make([]byte, len(data)),
	}
	copy(p.Data, data)
	return p
}

func checksum(data []byte) uint16 {
	var (
		sum    uint32
		length = len(data)
		index  int
	)
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length > 0 {
		sum += uint32(data[index])
	}
	sum += sum >> 16
	return uint16(^sum)
}
