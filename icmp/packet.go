package icmp

import (
	"bytes"
	"encoding/binary"
)

type PacketType byte

const (
	EchoRequest PacketType = 8
)

type Packet struct {
	Type       PacketType
	Code       byte
	Checksum   uint16
	Identifier uint16
	Data       []byte
}

func NewRequest() *Packet {
	p := &Packet{
		Type:       EchoRequest,
		Code:       0,
		Checksum:   0,
		Identifier: 0,
		Data:       nil,
	}
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, p)
	if err != nil {
		return nil
	}
	p.Checksum = checksum(buffer.Bytes())
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
