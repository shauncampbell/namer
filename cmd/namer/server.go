package namer

import (
	"fmt"
	"github.com/shauncampbell/namer/pkg/dns"
	"github.com/spf13/cobra"
)

var (
	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Start the DNS server",
		Long:  "Start the DNS server on the specified port",
		Run: func(cmd *cobra.Command, args []string) {
			namer := dns.NewServer(cfgFile, serverPort)
			if err := namer.Listen(); err != nil {
				fmt.Println(err.Error())
			}
		},
	}
	serverPort int
)
