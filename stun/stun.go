package stun

import (
	"fmt"
	"net"
	"net/netip"
)

type Opt struct {
	Server  netip.AddrPort
	IPQuery bool
}

func Run(o Opt) error {
	udp, err := net.ListenPacket("udp", "")
	if err != nil {
		return err
	}

	action := NoAction
	if o.IPQuery {

	}

	_, err = udp.WriteTo(MarshalRequest(NewRequest(action)), net.UDPAddrFromAddrPort(o.Server))
	if err != nil {
		return err
	}

	buff := make([]byte, 1460)
	n, addr, err := udp.ReadFrom(buff)
	if err != nil {
		return err
	}

	resp := &Response{}
	_, err = UnmarshalResponse(buff[:n], resp)
	if err != nil {
		return err
	}

	fmt.Printf("Recv STUN response from: %v\n", addr)
	ret := resp.Attributes[MappedAddress]
	fmt.Printf("External Address: %v:%d\n", ret.IP, ret.Port)
	return nil
}
