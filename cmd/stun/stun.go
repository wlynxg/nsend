package stun

import (
	"log"
	"net/netip"

	"github.com/spf13/cobra"
	"github.com/wlynxg/nsend/stun"
)

// Cmd represents the dns command
var Cmd = &cobra.Command{
	Use:   "stun",
	Short: "send STUN request to network hosts",
	Long:  `stun can send STUN request to the specified host。`,
	Run: func(cmd *cobra.Command, args []string) {
		dst, err := netip.ParseAddrPort(args[0])
		if err != nil {
			log.Fatalf("invalid addr: %s %v", args[0], err)
		}

		err = stun.Run(stun.Opt{
			Server: dst,
		})
		if err != nil {
			log.Fatalln(err)
		}
	},
	Args: cobra.ExactArgs(1),
}

func init() {
}
