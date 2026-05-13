package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var stopJobID string

// stopCmd represents the stop command.
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a job",
	Run: func(cmd *cobra.Command, args []string) {
		_, _ = apiClient.StopJob(stopJobID)
	},
}

func init() {
	clientCmd.AddCommand(stopCmd)
	stopCmd.Flags().StringVar(&stopJobID, "id", "", "ID of the job")
	if err := stopCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("could not prepare CLI flag, reason: %v", err)
	}
}
