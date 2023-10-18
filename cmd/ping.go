package cmd

import (
	"log"
	"math"
	"net"

	"github.com/spf13/cobra"
	"github.com/wlynxg/nsend/icmp"
)

var (
	PingCount int
)

// pingCmd represents the dns command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "send ICMP ECHO_REQUEST to network hosts",
	Long:  `ping uses the ICMP protocol's mandatory ECHO_REQUEST datagram to elicit an ICMP ECHO_RESPONSE from a host or gateway.`,
	Run: func(cmd *cobra.Command, args []string) {
		dst := net.ParseIP(args[0])
		if dst == nil {
			log.Fatalf("invalid addr: %s", args[0])
		}

		err := icmp.Run(icmp.Opt{
			Dst:   &net.IPAddr{IP: dst},
			Count: PingCount,
		})
		if err != nil {
			log.Fatalln(err)
		}
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	Root.AddCommand(pingCmd)
	pingCmd.Flags().IntVarP(&PingCount, "count", "c", math.MaxInt, "stop after <count> replies")
}
