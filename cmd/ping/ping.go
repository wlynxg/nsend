package ping

import (
	"log"
	"math"
	"net"

	"github.com/spf13/cobra"
	"github.com/wlynxg/nsend/icmp"
)

var (
	Count int
)

// Cmd represents the dns command
var Cmd = &cobra.Command{
	Use:   "ping [host]",
	Short: "send ICMP ECHO_REQUEST to network hosts",
	Long:  `ping uses the ICMP protocol's mandatory ECHO_REQUEST datagram to elicit an ICMP ECHO_RESPONSE from a host or gateway.`,
	Run: func(cmd *cobra.Command, args []string) {
		dst, err := net.ResolveIPAddr("ip", args[0])
		if err != nil {
			log.Fatalf("invalid addr: %s %v", args[0], err)
		}

		err = icmp.Run(icmp.Opt{
			DstRaw: args[0],
			Dst:    net.IPAddr{IP: dst.IP},
			Count:  Count,
		})
		if err != nil {
			log.Fatalln(err)
		}
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	Cmd.Flags().IntVarP(&Count, "count", "c", math.MaxInt, "stop after <count> replies")
}
