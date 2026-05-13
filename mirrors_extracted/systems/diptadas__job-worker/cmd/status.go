package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var queryJobID string

// statusCmd represents the get command.
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get status of a job",
	Run: func(cmd *cobra.Command, args []string) {
		_, _ = apiClient.GetJobStatus(queryJobID)
	},
}

func init() {
	clientCmd.AddCommand(statusCmd)
	statusCmd.Flags().StringVar(&queryJobID, "id", "", "ID of the job")
	if err := statusCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("could not prepare CLI flag, reason: %v\n", err)
	}
}
