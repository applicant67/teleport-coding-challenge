package cmd

import (
	"context"
	"errors"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/xSudoNymx/job-manager/gen/jobpb"
)

var streamCmd = &cobra.Command{
	Use:     "stream <job-id>",
	Short:   "Stream job output",
	Long:    "Stream stdout and stderr from a job in real-time.",
	Example: `  jobctl stream 550e8400-e29b-41d4-a716-446655440000`,
	Args:    cobra.ExactArgs(1),
	RunE:    runStream,
}

func init() {
	rootCmd.AddCommand(streamCmd)
}

func runStream(cmd *cobra.Command, args []string) error {
	jobID := args[0]

	client, err := newClient()
	if err != nil {
		return err
	}
	defer client.Close()

	// Use background context with signal handling (no timeout for streaming)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	stream, err := client.StreamOutput(ctx, jobID)
	if err != nil {
		return err
	}

	for {
		output, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, context.Canceled) {
				return nil
			}
			return err
		}

		switch output.Tag {
		case jobpb.StreamTag_STREAM_TAG_STDOUT:
			os.Stdout.Write(output.Data)
		case jobpb.StreamTag_STREAM_TAG_STDERR:
			os.Stderr.Write(output.Data)
		}
	}
}
