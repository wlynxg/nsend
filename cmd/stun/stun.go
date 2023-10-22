package stun

import (
	"fmt"
	"net"
	"net/netip"
	"time"

	"github.com/wlynxg/nsend/stun"
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

	action := stun.NoAction
	_, err = udp.WriteTo(stun.MarshalRequest(stun.NewRequest(action)), net.UDPAddrFromAddrPort(o.Server))
	if err != nil {
		return err
	}

	buff := make([]byte, 1460)
	n, addr, err := udp.ReadFrom(buff)
	if err != nil {
		return err
	}

	resp := &stun.Response{}
	_, err = stun.UnmarshalResponse(buff[:n], resp)
	if err != nil {
		return err
	}

	fmt.Printf("Recv STUN response from: %v\n", addr)
	ret := resp.Attributes[stun.MappedAddress]
	fmt.Printf("External Address: %v:%d\n", ret.IP, ret.Port)
	return nil
}

func natTypeDetect(udp *net.UDPConn, dst *net.UDPAddr) error {
	var (
		buff                       = make([]byte, 1460)
		resp                       = &stun.Response{}
		other, addr1, addr2, addr3 netip.AddrPort
		filter                     stun.FilteringBehavior
		mapped                     stun.MappingBehavior
	)

	// Step1
	_, err := udp.WriteTo(stun.MarshalRequest(stun.NewRequest(stun.NoAction)), dst)
	if err != nil {
		return err
	}

	udp.SetReadDeadline(time.Now().Add(timeout))
	n, err := udp.Read(buff)
	if err != nil {
		return err
	}

	_, err = stun.UnmarshalResponse(buff[:n], resp)
	if err != nil {
		return err
	}

	if v, ok := resp.Attributes[stun.MappedAddress]; ok {
		addr1 = netip.AddrPortFrom(v.IP, uint16(v.Port))
	}

	if v, ok := resp.Attributes[stun.OtherAddress]; ok {
		other = netip.AddrPortFrom(v.IP, uint16(v.Port))
	} else if v, ok = resp.Attributes[stun.ChangedAddress]; ok {
		other = netip.AddrPortFrom(v.IP, uint16(v.Port))
	}

	// Step2
	_, err = udp.WriteTo(stun.MarshalRequest(stun.NewRequest(stun.ChangeIPAndPort)), dst)
	if err != nil {
		return err
	}

	udp.SetReadDeadline(time.Now().Add(timeout))
	n, err = udp.Read(buff)
	if err == nil {
		filter = stun.EndpointIndependentFiltering
		goto step4
	}

	// Step3
	_, err = udp.WriteTo(stun.MarshalRequest(stun.NewRequest(stun.ChangePort)), dst)
	if err != nil {
		return err
	}

	udp.SetReadDeadline(time.Now().Add(timeout))
	n, err = udp.Read(buff)
	if err != nil {
		filter = stun.AddressAndPortDependentFiltering
	} else {
		filter = stun.AddressDependentFiltering
	}

step4:
	if addr1 == udp.LocalAddr().(*net.UDPAddr).AddrPort() {
		mapped = stun.NoMapping
		goto ret
	}

	// Step4
	_, err = udp.WriteTo(stun.MarshalRequest(stun.NewRequest(stun.NoAction)),
		net.UDPAddrFromAddrPort(netip.AddrPortFrom(other.Addr(), uint16(dst.Port))))
	if err != nil {
		return err
	}

	udp.SetReadDeadline(time.Now().Add(timeout))
	n, err = udp.Read(buff)
	if err != nil {
		return err
	}

	if v, ok := resp.Attributes[stun.MappedAddress]; ok {
		addr2 = netip.AddrPortFrom(v.IP, uint16(v.Port))
	}

	if addr1 == addr2 {
		mapped = stun.EndpointIndependentMapping
	}

	// Step5
	_, err = udp.WriteTo(stun.MarshalRequest(stun.NewRequest(stun.NoAction)), net.UDPAddrFromAddrPort(other))
	if err != nil {
		return err
	}

	udp.SetReadDeadline(time.Now().Add(timeout))
	n, err = udp.Read(buff)
	if err != nil {
		return err
	}

	if v, ok := resp.Attributes[stun.MappedAddress]; ok {
		addr3 = netip.AddrPortFrom(v.IP, uint16(v.Port))
	}

	if addr2 == addr3 {
		mapped = stun.AddressDependentMapping
	} else {
		mapped = stun.AddressAndPortDependentMapping
	}

ret:

	fmt.Println(mapped)
	fmt.Println(filter)
	return nil
}
