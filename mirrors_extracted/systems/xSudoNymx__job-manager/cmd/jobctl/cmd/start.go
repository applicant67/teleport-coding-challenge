package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [flags] -- <command> [args...]",
	Short: "Start a new job",
	Long:  "Start a new job with the specified command and arguments.",
	Example: `  jobctl start -- /bin/echo hello world
  jobctl start -- /bin/sleep 60`,
	Args: cobra.MinimumNArgs(1),
	RunE: runStart,
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func runStart(cmd *cobra.Command, args []string) error {
	exe := args[0]
	var cmdArgs []string
	if len(args) > 1 {
		cmdArgs = args[1:]
	}

	client, err := newClient()
	if err != nil {
		return err
	}
	defer client.Close()

	ctx, cancel := newContext()
	defer cancel()

	resp, err := client.Start(ctx, exe, cmdArgs)
	if err != nil {
		return fmt.Errorf("start job: %w", err)
	}

	fmt.Println(resp.JobId)
	return nil
}
