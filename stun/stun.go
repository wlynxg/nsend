package stun

import (
	"fmt"
	"net"
	"net/netip"
	"time"
)

var (
	timeout = 5 * time.Second
)

type Opt struct {
	Server        netip.AddrPort
	NatTypeDetect bool
}

func Run(o Opt) error {
	udp, err := net.ListenUDP("udp", nil)
	if err != nil {
		return err
	}

	if o.NatTypeDetect {
		return natTypeDetect(udp, net.UDPAddrFromAddrPort(o.Server))
	}

	action := NoAction
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

func natTypeDetect(udp *net.UDPConn, dst *net.UDPAddr) error {
	var (
		buff                       = make([]byte, 1460)
		resp                       = &Response{}
		other, addr1, addr2, addr3 netip.AddrPort
		filter                     FilteringBehavior
		mapped                     MappingBehavior
	)

	// Step1
	_, err := udp.WriteTo(MarshalRequest(NewRequest(NoAction)), dst)
	if err != nil {
		return err
	}

	udp.SetReadDeadline(time.Now().Add(timeout))
	n, err := udp.Read(buff)
	if err != nil {
		return err
	}

	_, err = UnmarshalResponse(buff[:n], resp)
	if err != nil {
		return err
	}

	if v, ok := resp.Attributes[MappedAddress]; ok {
		addr1 = netip.AddrPortFrom(v.IP, uint16(v.Port))
	}

	if v, ok := resp.Attributes[OtherAddress]; ok {
		other = netip.AddrPortFrom(v.IP, uint16(v.Port))
	}

	// Step2
	_, err = udp.WriteTo(MarshalRequest(NewRequest(ChangeIPAndPort)), dst)
	if err != nil {
		return err
	}

	udp.SetReadDeadline(time.Now().Add(timeout))
	n, err = udp.Read(buff)
	if err == nil {
		filter = EndpointIndependentFiltering
		goto step4
	}

	// Step3
	_, err = udp.WriteTo(MarshalRequest(NewRequest(ChangePort)), dst)
	if err != nil {
		return err
	}

	udp.SetReadDeadline(time.Now().Add(timeout))
	n, err = udp.Read(buff)
	if err != nil {
		filter = AddressAndPortDependentFiltering
	} else {
		filter = AddressDependentFiltering
	}

step4:
	if addr1 == udp.LocalAddr().(*net.UDPAddr).AddrPort() {
		mapped = NoMapping
		goto ret
	}

	// Step4
	_, err = udp.WriteTo(MarshalRequest(NewRequest(NoAction)),
		net.UDPAddrFromAddrPort(netip.AddrPortFrom(other.Addr(), uint16(dst.Port))))
	if err != nil {
		return err
	}

	udp.SetReadDeadline(time.Now().Add(timeout))
	n, err = udp.Read(buff)
	if err != nil {
		return err
	}

	if v, ok := resp.Attributes[MappedAddress]; ok {
		addr2 = netip.AddrPortFrom(v.IP, uint16(v.Port))
	}

	if addr1 == addr2 {
		mapped = EndpointIndependentMapping
	}

	// Step5
	_, err = udp.WriteTo(MarshalRequest(NewRequest(NoAction)), net.UDPAddrFromAddrPort(other))
	if err != nil {
		return err
	}

	udp.SetReadDeadline(time.Now().Add(timeout))
	n, err = udp.Read(buff)
	if err != nil {
		return err
	}

	if v, ok := resp.Attributes[MappedAddress]; ok {
		addr3 = netip.AddrPortFrom(v.IP, uint16(v.Port))
	}

	if addr2 == addr3 {
		mapped = AddressDependentMapping
	} else {
		mapped = AddressAndPortDependentMapping
	}

ret:

	fmt.Println(mapped)
	fmt.Println(filter)
	return nil
}
