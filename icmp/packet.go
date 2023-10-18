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

type Packet struct {
	Type       PacketType
	Code       byte
	Checksum   uint16
	Identifier uint16
	Data       []byte
}

func MarshalPacket(p *Packet) []byte {
	if p == nil {
		return nil
	}

	buff := make([]byte, 512)
	buff[0] = byte(p.Type)
	buff[1] = p.Code
	binary.BigEndian.PutUint16(buff[2:], p.Checksum)
	binary.BigEndian.PutUint16(buff[4:], p.Identifier)
	copy(buff[6:], p.Data)
	return buff[6+len(p.Data):]
}

func UnmarshalPacket(buff []byte, p *Packet) (int, error) {
	if p == nil {
		return -1, errors.New("response cannot be nil")
	}

	if len(buff) < 6 {
		return -1, errors.New("not a correct icmp package")
	}

	p.Type = PacketType(buff[0])
	p.Code = buff[1]
	p.Checksum = binary.BigEndian.Uint16(buff[2:4])
	p.Identifier = binary.BigEndian.Uint16(buff[4:6])
	if len(buff) > 6 {
		data := make([]byte, len(buff)-6)
		copy(data, buff[6:])
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
		Data:       nil,
	}

	buff := make([]byte, 512)
	buff[0] = byte(p.Type)
	buff[1] = p.Code
	binary.BigEndian.PutUint16(buff[2:], p.Checksum)
	binary.BigEndian.PutUint16(buff[4:], p.Identifier)
	copy(buff[6:], data)
	p.Checksum = checksum(buff[6+len(data):])
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
