package cmd

import (
	"job-worker/api/server"

	"github.com/spf13/cobra"
)

var apiServer server.Server

// serverCmd represents the server command.
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the job-worker API server",
	Run: func(cmd *cobra.Command, args []string) {
		apiServer.InitializeAndRun()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntVarP(&apiServer.Port, "port", "p", 8443, "Port of API server")
	serverCmd.Flags().StringVar(&apiServer.CaCert, "ca", "ssl/ca.crt", "CA cert file")
	serverCmd.Flags().StringVar(&apiServer.ServerCert, "cert", "ssl/server.crt", "Server cert file")
	serverCmd.Flags().StringVar(&apiServer.ServerKey, "key", "ssl/server.key", "Server key file")
}
