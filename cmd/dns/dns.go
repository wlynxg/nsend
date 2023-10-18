package dns

import (
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/wlynxg/nsend/dns"
)

// Cmd represents the dns command
var Cmd = &cobra.Command{
	Use:   "dns",
	Short: "dns domain name query",
	Long:  `Query the resolution address of the domain name, you can specify the type and dns server`,
	Run: func(cmd *cobra.Command, args []string) {
		t := dns.A
		switch cmd.Flags().Lookup("type").Value.String() {
		case "A":
			t = dns.A
		}

		response, err := dns.Run(dns.Opt{
			Domain: args[0],
			Type:   t,
			Server: cmd.Flags().Lookup("server").Value.String(),
		})
		if err != nil {
			panic(err)
		}

		// Conditional Time column.
		table := tablewriter.NewWriter(color.Output)
		table.SetHeader([]string{"Name", "Type", "Class", "TTL", "Address"})
		table.SetAutoWrapText(true)
		table.SetAutoFormatHeaders(true)
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetCenterSeparator("")
		table.SetColumnSeparator("")
		table.SetRowSeparator("")
		table.SetHeaderLine(false)
		table.SetBorder(false)
		table.SetTablePadding("\t") // pad with tabs
		table.SetNoWhiteSpace(true)

		var (
			red   = color.New(color.FgRed, color.Bold).SprintFunc()
			green = color.New(color.FgGreen, color.Bold).SprintFunc()
		)

		for _, answer := range response.Answers {
			table.Append([]string{green(answer.Name), answer.Type.String(), answer.Class.String(),
				cast.ToString(answer.TTL), red(answer.Address.String())})
		}
		table.Render()
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	Cmd.Flags().StringP("server", "s", "", "Specify dns server, format: protocol:[address]:port")
	Cmd.Flags().StringP("type", "t", "A", "Specify the domain name record type, the default is A record")
}
