package stun

import (
	"log"
	"net"
	"net/netip"

	"github.com/spf13/cobra"
	"github.com/wlynxg/nsend/stun"
)

var (
	port uint16
)

// Cmd represents the dns command
var Cmd = &cobra.Command{
	Use:   "stun <STUN server>",
	Short: "send STUN request to network hosts",
	Long:  `stun can send STUN request to the specified host。`,
	Run: func(cmd *cobra.Command, args []string) {
		dst, err := net.ResolveIPAddr("ip", args[0])
		if err != nil {
			log.Fatalf("invalid addr: %s %v", args[0], err)
		}

		err = stun.Run(stun.Opt{
			Server: netip.AddrPortFrom(netip.MustParseAddr(dst.String()), port),
		})
		if err != nil {
			log.Fatalln(err)
		}
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	Cmd.Flags().Uint16VarP(&port, "port", "p", stun.DefaultSTUNPort, "stun server port")
}
