package wol

import (
	"net"

	"github.com/pkg/errors"
)

type Opt struct {
	DstMac   net.HardwareAddr
	Password []byte

	IP        net.IP
	Port      int
	Interface *net.Interface
}

const (
	DefaultPort = 9
)

func Run(o Opt) error {
	if len(o.DstMac) == 0 {
		return errors.New("DstMac cannot be empty")
	}

	switch len(o.Password) {
	case 0, 4, 6:
	default:
		return errors.New("The password length must be 0, 4, 6")
	}

	if o.Interface != nil {
		//	TODO: use raw socket
	}

	if o.IP == nil {
		o.IP = net.IPv4bcast
	}

	if o.Port == 0 {
		o.Port = DefaultPort
	}

	udp, err := net.DialUDP("udp", &net.UDPAddr{IP: o.IP, Port: DefaultPort}, nil)
	if err != nil {
		return err
	}

	_, err = udp.Write(MarshalRequest(NewRequest(o.DstMac, o.Password)))
	if err != nil {
		return err
	}

	return nil
}
