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
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
