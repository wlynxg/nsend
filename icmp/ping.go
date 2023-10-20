package icmp

import (
	"fmt"
	"log"
	"net"
	"time"
)

type Opt struct {
	DstRaw string
	Dst    net.IPAddr
	Count  int
	Data   []byte
}

func Run(o Opt) error {
	client, err := net.DialIP("ip4:icmp", nil, &o.Dst)
	if err != nil {
		return err
	}

	var (
		req  = NewRequest(o.Data)
		buff = make([]byte, 1460)

		resp     Packet
		start    time.Time
		duration time.Duration
	)

	fmt.Printf("PING %s (%s) %d bytes of data.", o.DstRaw, o.Dst.String(), len(o.Data))
	for i := 0; i < o.Count; i++ {
		start = time.Now()

		req.Sequence++
		_, err = client.Write(MarshalPacket(req))
		if err != nil {
			return err
		}

		err = client.SetReadDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			log.Fatal(err)
		}

		n, addr, err := client.ReadFrom(buff)
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
				fmt.Printf("from %s (%v): icmp_seq=%d time=%v\t\n", o.DstRaw, addr.String(), req.Sequence, duration)
			}
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}
