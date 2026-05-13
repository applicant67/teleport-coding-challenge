package cmd

import (
	"job-worker/api/client"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var apiClient client.Client

// clientCmd represents the client command.
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Client for job-worker service",
}

func init() {
	rootCmd.AddCommand(clientCmd)
	clientCmd.PersistentFlags().StringVar(&apiClient.Address, "address", "https://localhost:8443", "Address of the API server")
	clientCmd.PersistentFlags().StringVar(&apiClient.CaCert, "ca", "ssl/ca.crt", "CA cert file")
	clientCmd.PersistentFlags().StringVar(&apiClient.ClientCert, "cert", "", "Client cert file")
	clientCmd.PersistentFlags().StringVar(&apiClient.ClientKey, "key", "", "Client key file")

	if err := clientCmd.MarkPersistentFlagRequired("cert"); err != nil {
		log.Fatalf("could not prepare CLI flag, reason: %v\n", err)
	}
	if err := clientCmd.MarkPersistentFlagRequired("key"); err != nil {
		log.Fatalf("could not prepare CLI flag, reason: %v\n", err)
	}
}
