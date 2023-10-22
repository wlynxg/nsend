package stun

import (
	"bytes"
	"fmt"
	"log"
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
		resp                       = &stun.Response{}
		req                        = &stun.Request{}
		other, addr1, addr2, addr3 netip.AddrPort
		filter                     stun.FilteringBehavior
		mapped                     stun.MappingBehavior
	)

	// Step1
	req = stun.NewRequest(stun.NoAction)
	_, err := udp.WriteTo(stun.MarshalRequest(req), dst)
	if err != nil {
		return err
	}

	success, err := waitResponse(udp, req.TransactionID[:], resp)
	if err != nil {
		return err
	}

	if !success {
		log.Println(stun.UDPBlock)
		return nil
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
	req = stun.NewRequest(stun.ChangeIPAndPort)
	_, err = udp.WriteTo(stun.MarshalRequest(req), dst)
	if err != nil {
		return err
	}

	success, err = waitResponse(udp, req.TransactionID[:], resp)
	if err != nil {
		return err
	}

	if success {
		filter = stun.EndpointIndependentFiltering
		goto step4
	}

	// Step3
	req = stun.NewRequest(stun.ChangePort)
	_, err = udp.WriteTo(stun.MarshalRequest(req), dst)
	if err != nil {
		return err
	}

	success, err = waitResponse(udp, req.TransactionID[:], resp)
	if err != nil {
		return err
	}

	if success {
		filter = stun.AddressDependentFiltering
	} else {
		filter = stun.AddressAndPortDependentFiltering
	}

step4:
	if addr1 == udp.LocalAddr().(*net.UDPAddr).AddrPort() {
		mapped = stun.NoMapping
		goto ret
	}

	// Step4
	req = stun.NewRequest(stun.NoAction)
	_, err = udp.WriteTo(stun.MarshalRequest(req),
		net.UDPAddrFromAddrPort(netip.AddrPortFrom(other.Addr(), uint16(dst.Port))))
	if err != nil {
		return err
	}

	success, err = waitResponse(udp, req.TransactionID[:], resp)
	if err != nil {
		return err
	}

	if success {
		log.Println(stun.UDPBlock)
		return nil
	}

	if v, ok := resp.Attributes[stun.MappedAddress]; ok {
		addr2 = netip.AddrPortFrom(v.IP, uint16(v.Port))
	}

	if addr1 == addr2 {
		mapped = stun.EndpointIndependentMapping
	}

	// Step5
	req = stun.NewRequest(stun.NoAction)
	_, err = udp.WriteTo(stun.MarshalRequest(req), net.UDPAddrFromAddrPort(other))
	if err != nil {
		return err
	}

	success, err = waitResponse(udp, req.TransactionID[:], resp)
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

// success, err
func waitResponse(udp *net.UDPConn, txid stun.TxID, resp *stun.Response) (bool, error) {
	buff := make([]byte, 1460)
	for {
		udp.SetReadDeadline(time.Now().Add(timeout))
		n, err := udp.Read(buff)
		if err != nil {
			netErr, ok := err.(*net.OpError)
			if ok && netErr.Timeout() {
				return false, nil
			}
			return false, err
		}

		_, err = stun.UnmarshalResponse(buff[:n], resp)
		if err != nil {
			return false, err
		}

		if bytes.Equal(resp.TransactionID, txid) {
			return true, nil
		}
	}
}
