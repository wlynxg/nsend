package icmp

import (
	"fmt"
	"log"
	"net"
	"time"
)

type Opt struct {
	Dst   net.Addr
	Count int
	Data  []byte
}

func Run(o Opt) error {
	listen, err := net.ListenPacket("ip:icmp", "")
	if err != nil {
		return err
	}

	var (
		req      = NewRequest(o.Data)
		resp     Packet
		buff     = make([]byte, 1500)
		start    time.Time
		duration time.Duration
	)

	for i := 0; i < o.Count; i++ {
		start = time.Now()
		_, err = listen.WriteTo(MarshalPacket(req), o.Dst)
		if err != nil {
			return err
		}

		err = listen.SetReadDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			log.Fatal(err)
		}

		n, addr, err := listen.ReadFrom(buff)
		if err != nil {
			return err
		}
		duration = time.Since(start)

		_, err = UnmarshalPacket(buff[:n], &resp)
		if err != nil {
			return err
		}

		switch resp.Type {
		case EchoReply:
			if addr.String() == o.Dst.String() {
				fmt.Printf("reply from %s: time=%v\n", addr.String(), duration)
			}
		}
	}
	return nil
}
