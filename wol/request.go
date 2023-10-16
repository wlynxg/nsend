package wol

import (
	"net"
)

type Request struct {
	Header   [12]byte
	Dst      net.HardwareAddr
	Password []byte
}

func MarshalRequest(req *Request) []byte {
	if req == nil {
		return nil
	}

	var (
		buff   = make([]byte, 512)
		offset = 0
	)

	copy(buff[:], req.Header[:])
	offset += 12
	for i := 0; i < 6; i++ {
		copy(buff[offset:], req.Dst)
		offset += len(req.Dst)
	}

	copy(buff[offset:], req.Password)
	offset += len(req.Password)
	return buff[:offset]
}

func NewRequest(dst net.HardwareAddr, pwd []byte) *Request {
	switch len(pwd) {
	case 0, 4, 6:
	default:
		return nil
	}

	req := &Request{
		Header:   [12]byte{},
		Dst:      dst,
		Password: pwd,
	}
	return req
}
