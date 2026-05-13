package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:     "stop <job-id>",
	Short:   "Stop a running job",
	Long:    "Stop a running job by its ID.",
	Example: `  jobctl stop 550e8400-e29b-41d4-a716-446655440000`,
	Args:    cobra.ExactArgs(1),
	RunE:    runStop,
}

func init() {
	rootCmd.AddCommand(stopCmd)
}

func runStop(cmd *cobra.Command, args []string) error {
	jobID := args[0]

	client, err := newClient()
	if err != nil {
		return err
	}
	defer client.Close()

	ctx, cancel := newContext()
	defer cancel()

	_, err = client.Stop(ctx, jobID)
	if err != nil {
		return fmt.Errorf("stop job: %w", err)
	}

	fmt.Println("stopped")
	return nil
}
