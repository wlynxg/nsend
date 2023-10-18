package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/wlynxg/nsend/cmd/dns"
	"github.com/wlynxg/nsend/cmd/ping"
	"github.com/wlynxg/nsend/cmd/wol"
)

// Root represents the base command when called without any subcommands
var Root = &cobra.Command{
	Use:   "nsend",
	Short: "Versatile command-line tool for sending various network packets.",
	Long: `This command line tool is a feature-rich network packet launcher 
that is capable of sending various types of network packets.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the Root.
func Execute() {
	err := Root.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	Root.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	Root.AddCommand(dns.Cmd)
	Root.AddCommand(wol.Cmd)
	Root.AddCommand(ping.Cmd)
}
