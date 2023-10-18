package cmd

import (
	"fmt"
	"net"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/wlynxg/nsend/wol"
)

// dnsCmd represents the dns command
var wolCmd = &cobra.Command{
	Use:   "wol",
	Short: "wol wakes up the LAN host",
	Long: `Send a wol request to wake up the host in the LAN. You can choose to use udp method or raw socket method.
When the interface is specified, the raw socket method is used; if not specified, the udp method is used.`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			dst  net.HardwareAddr
			pwd  []byte
			ifi  *net.Interface
			ip   net.IP
			port int
		)

		dst, err := net.ParseMAC(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}

		pwdStr := cmd.Flags().Lookup("password").Value.String()
		if pwdStr != "" {
			switch len(pwdStr) {
			case 0, 4, 6:
				pwd = []byte(pwdStr)
			default:
				fmt.Println(errors.New("The password length must be 0, 4, 6"))
				return
			}
		}

		ifiName := cmd.Flags().Lookup("interface").Value.String()
		if ifiName != "" {
			ifi, err = net.InterfaceByName(ifiName)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		ipStr := cmd.Flags().Lookup("ip").Value.String()
		if ipStr != "" {
			ip = net.ParseIP(ipStr)
		}

		portStr := cmd.Flags().Lookup("ip").Value.String()
		if portStr != "" {
			port, err = cast.ToIntE(portStr)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		err = wol.Run(wol.Opt{
			DstMac:    dst,
			Password:  pwd,
			IP:        ip,
			Port:      port,
			Interface: ifi,
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("sent Wake-on-LAN magic packet to %s", dst.String())
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	Root.AddCommand(wolCmd)
	wolCmd.Flags().StringP("interface", "i", "", "network interface to use to send Wake-on-LAN magic packet")
	wolCmd.Flags().StringP("password", "p", "", "optional password for Wake-on-LAN magic packet")
	wolCmd.Flags().String("ip", "p", "network address for Wake-on-LAN magic packet")
	wolCmd.Flags().Int("port", 9, "network port for Wake-on-LAN magic packet")
}
