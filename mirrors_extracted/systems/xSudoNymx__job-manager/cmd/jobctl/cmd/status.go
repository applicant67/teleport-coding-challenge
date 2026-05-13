package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/xSudoNymx/job-manager/gen/jobpb"
)

var statusCmd = &cobra.Command{
	Use:     "status <job-id>",
	Short:   "Get job status",
	Long:    "Get the current status of a job by its ID.",
	Example: `  jobctl status 550e8400-e29b-41d4-a716-446655440000`,
	Args:    cobra.ExactArgs(1),
	RunE:    runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) error {
	jobID := args[0]

	client, err := newClient()
	if err != nil {
		return err
	}
	defer client.Close()

	ctx, cancel := newContext()
	defer cancel()

	resp, err := client.Status(ctx, jobID)
	if err != nil {
		return fmt.Errorf("get status: %w", err)
	}

	printStatus(resp.Status)
	return nil
}

func printStatus(s *jobpb.JobStatus) {
	fmt.Printf("Job ID:     %s\n", s.JobId)
	fmt.Printf("PID:        %d\n", s.Pid)
	fmt.Printf("State:      %s\n", formatState(s.State))
	fmt.Printf("Started:    %s\n", formatTime(s.StartedAt))
	if s.EndedAt > 0 {
		fmt.Printf("Ended:      %s\n", formatTime(s.EndedAt))
	}
	if s.State == jobpb.JobState_JOB_STATE_EXITED || s.State == jobpb.JobState_JOB_STATE_STOPPED {
		fmt.Printf("Exit Code:  %d\n", s.ExitCode)
	}
	if s.WaitErr != "" {
		fmt.Printf("Wait Error: %s\n", s.WaitErr)
	}
}

func formatState(s jobpb.JobState) string {
	switch s {
	case jobpb.JobState_JOB_STATE_RUNNING:
		return "running"
	case jobpb.JobState_JOB_STATE_STOPPING:
		return "stopping"
	case jobpb.JobState_JOB_STATE_STOPPED:
		return "stopped"
	case jobpb.JobState_JOB_STATE_EXITED:
		return "exited"
	default:
		return "unknown"
	}
}

func formatTime(nanos int64) string {
	if nanos == 0 {
		return "-"
	}
	return time.Unix(0, nanos).Format(time.RFC3339)
}
