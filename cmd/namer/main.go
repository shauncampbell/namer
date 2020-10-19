package namer

import "github.com/spf13/cobra"

var (
	rootCmd = &cobra.Command{
		Use:   "dapper",
		Short: "Dapper is a very fast simple ldap server",
		Long:  `A fast easy to use LDAP server for use with home labs`,
	}

	cfgFile    string

)

func init() {
	cobra.OnInitialize()

	// Add flags to server command
	serverCmd.PersistentFlags().IntVarP(&serverPort, "port", "p", 53, "port to run DNS server on")
	// Add flags to the root command
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "f", "", "config file (default is $HOME/.namer.yaml)")
	rootCmd.MarkPersistentFlagRequired("config")

	// Add the sub commands to the root
	rootCmd.AddCommand(serverCmd)
}

func main() {
	_ = rootCmd.Execute()
}
