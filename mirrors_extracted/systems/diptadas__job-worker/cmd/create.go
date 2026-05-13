package cmd

import (
	"job-worker/lib"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var request lib.CreateJobRequest

// createCmd represents the create command.
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new job",
	Run: func(cmd *cobra.Command, args []string) {
		_, _ = apiClient.CreateJob(request)
	},
}

func init() {
	clientCmd.AddCommand(createCmd)
	createCmd.Flags().StringVar(&request.Command, "cmd", "", "Path to the executable file")
	createCmd.Flags().StringArrayVar(&request.Args, "args", nil, "Args to executable")

	if err := createCmd.MarkFlagRequired("cmd"); err != nil {
		log.Fatalf("could not prepare CLI flag, reason: %v\n", err)
	}
}
