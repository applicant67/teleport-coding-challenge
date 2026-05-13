// Package cli implements the command-line interface for the Job Worker Service.
package cli

import (
	"cmp"
	"context"
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/spf13/cobra"
	pb "github.com/teleport-handson/api/v1"
	"google.golang.org/grpc"
)

// Version represents the current version of the CLI.
// It is intended to be set at build time using -ldflags.
var Version = "dev"

// Config holds the global configuration flags and settings for the CLI.
type Config struct {
	// ServerAddr is the address of the gRPC server (e.g., "localhost:8080").
	ServerAddr string
	// CAPath is the path to the Certificate Authority file for verifying the server.
	CAPath string
	// CertPath is the path to the client's certificate file for mTLS.
	CertPath string
	// KeyPath is the path to the client's private key file for mTLS.
	KeyPath string
	// OutputJSON indicates whether the output should be formatted as JSON.
	OutputJSON bool
	// Quiet suppresses informational output, printing only IDs or raw data where applicable.
	Quiet bool
	// WorkDir is the working directory for the job (start command).
	WorkDir string
	// Follow indicates if logs should be streamed after starting (start command).
	Follow bool
	// Connect is the factory function for establishing the gRPC connection.
	// It is injectable to allow mocking in tests.
	Connect func() (pb.JobWorkerClient, *grpc.ClientConn, error)
}

const (
	formatStart  = "Job initiated successfully.\nID:      %s\nStatus:  %s\n"
	formatStop   = "Job ID '%s' stopped successfully.\n"
)

// NewRootCmd creates the root command for the worker-cli.
// It initializes the configuration, flags, and subcommands.
func NewRootCmd() *cobra.Command {
	cfg := &Config{}
	cfg.Connect = cfg.dial

	rootCmd := &cobra.Command{
		Use:   "worker-cli",
		Short: "CLI for the Job Worker Service",
		Long: `A command-line interface for interacting with the Job Worker Service.
It supports starting, stopping, querying status, and streaming logs of isolated Linux processes.
Authentication is enforced via mTLS.`,
		Example: `  # Start a job
  worker-cli start -- echo "hello world"

  # Stream logs for a specific job
  worker-cli stream <job_id>

  # Check status
  worker-cli status <job_id>

  # Stop a job
  worker-cli stop <job_id>`,
		// SilenceUsage prevents Cobra from printing the help message when a command returns an error.
		// This ensures that runtime errors (e.g., connection refused) don't look like syntax errors.
		SilenceUsage: true,
	}

	// Global flags for connection and mTLS
	rootCmd.PersistentFlags().StringVar(&cfg.ServerAddr, flagNameServerAddr, cmp.Or(os.Getenv(envVarServerAddr), defaultServerAddr), helpServerAddr)
	rootCmd.PersistentFlags().StringVar(&cfg.CAPath, flagNameCAPath, cmp.Or(os.Getenv(envVarCAPath), defaultCAPath), helpCAPath)
	rootCmd.PersistentFlags().StringVar(&cfg.CertPath, flagNameCertPath, cmp.Or(os.Getenv(envVarCertPath), defaultCertPath), helpCertPath)
	rootCmd.PersistentFlags().StringVar(&cfg.KeyPath, flagNameKeyPath, cmp.Or(os.Getenv(envVarKeyPath), defaultKeyPath), helpKeyPath)

	// Output formatting
	rootCmd.PersistentFlags().BoolVar(&cfg.OutputJSON, flagNameJSON, false, helpJSON)
	rootCmd.PersistentFlags().BoolVarP(&cfg.Quiet, flagNameQuiet, "q", false, helpQuiet)

	rootCmd.AddCommand(newStartCommand(cfg))
	rootCmd.AddCommand(newStopCommand(cfg))
	rootCmd.AddCommand(newStatusCommand(cfg))
	rootCmd.AddCommand(newStreamCommand(cfg))
	rootCmd.AddCommand(newVersionCommand())

	return rootCmd
}

// newStartCommand creates the 'start' command to initiate a new job.
func newStartCommand(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start -- <command> [args...]",
		Short: "Start a new job",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cfg.withClient(func(client pb.JobWorkerClient) error {
				return cfg.startJob(cmd.Context(), cmd.OutOrStdout(), client, args[0], args[1:])
			})
		},
	}
	cmd.Flags().StringVar(&cfg.WorkDir, flagNameWorkDir, defaultWorkDir, helpWorkDir)
	cmd.Flags().BoolVarP(&cfg.Follow, flagNameFollow, "f", false, helpFollow)
	return cmd
}

// newStopCommand creates the 'stop' command to terminate a running job.
func newStopCommand(cfg *Config) *cobra.Command {
	return &cobra.Command{
		Use:   "stop <job_id>",
		Short: "Stop a running job",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cfg.withClient(func(client pb.JobWorkerClient) error {
				return cfg.stopJob(cmd.Context(), cmd.OutOrStdout(), client, args[0])
			})
		},
	}
}

// newStatusCommand creates the 'status' command to retrieve job details.
func newStatusCommand(cfg *Config) *cobra.Command {
	return &cobra.Command{
		Use:   "status <job_id>",
		Short: "Get job status",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cfg.withClient(func(client pb.JobWorkerClient) error {
				return cfg.getJobStatus(cmd.Context(), cmd.OutOrStdout(), client, args[0])
			})
		},
	}
}

// newStreamCommand creates the 'stream' command to stream job logs.
func newStreamCommand(cfg *Config) *cobra.Command {
	return &cobra.Command{
		Use:   "stream <job_id>",
		Short: "Stream job logs",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cfg.withClient(func(client pb.JobWorkerClient) error {
				return cfg.streamLogs(cmd.Context(), cmd.OutOrStdout(), client, args[0])
			})
		},
	}
}

// newVersionCommand creates the 'version' command to display CLI version.
func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("worker-cli version: %s\n", Version)
		},
	}
}

// streamLogs streams output from a job to the provided writer.
func (c *Config) streamLogs(ctx context.Context, out io.Writer, client pb.JobWorkerClient, jobID string) error {
	stream, err := client.Stream(ctx, &pb.StreamRequest{JobId: jobID})
	if err != nil {
		return fmt.Errorf("stream failed: %w", err)
	}

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			return fmt.Errorf("stream error: %w", err)
		}
		if _, err := out.Write(chunk.Data); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
	}
}

// startJob handles the logic for the 'start' command.
func (c *Config) startJob(ctx context.Context, out io.Writer, client pb.JobWorkerClient, command string, args []string) error {
	req := &pb.StartRequest{
		Command:          command,
		Args:             args,
		WorkingDirectory: c.WorkDir,
	}

	resp, err := client.Start(ctx, req)
	if err != nil {
		return fmt.Errorf("start failed: %w", err)
	}

	if c.OutputJSON {
		return printJSON(out, resp)
	}

	if c.Quiet {
		fmt.Fprintln(out, resp.JobId)
	} else {
		fmt.Fprintf(out, formatStart, resp.JobId, resp.Status)
	}
	if c.Follow {
		return c.streamLogs(ctx, out, client, resp.JobId)
	}
	return nil
}

// stopJob handles the logic for the 'stop' command.
func (c *Config) stopJob(ctx context.Context, out io.Writer, client pb.JobWorkerClient, jobID string) error {
	resp, err := client.Stop(ctx, &pb.StopRequest{JobId: jobID})
	if err != nil {
		return fmt.Errorf("stop failed: %w", err)
	}

	if c.OutputJSON {
		return printJSON(out, resp)
	}

	if c.Quiet {
		fmt.Fprintln(out, jobID)
	} else {
		fmt.Fprintf(out, formatStop, jobID)
	}
	return nil
}

const statusTemplate = `
Job Information:
  ID:            {{.ID}}
  Owner:         {{.Owner}}
  Command:       {{.Command}}
  Args:          {{.Args}}

Lifecycle:
  Status:        {{.Status}}
  Started:       {{.Started}}
  Ended:         {{.Ended}}
  Exit Code:     {{.ExitCode}}

Resource Constraints:
  CPU Limit:     {{printf "%.2f" .CPULimit}}%
  Memory Limit:  {{.MemoryLimit}} bytes
  Disk I/O:      Read: {{.DiskReadBPS}} bps | Write: {{.DiskWriteBPS}} bps
`

// getJobStatus handles the logic for the 'status' command.
func (c *Config) getJobStatus(ctx context.Context, out io.Writer, client pb.JobWorkerClient, jobID string) error {
	resp, err := client.Status(ctx, &pb.StatusRequest{JobId: jobID})
	if err != nil {
		return fmt.Errorf("status failed: %w", err)
	}

	if c.OutputJSON {
		return printJSON(out, resp)
	}

	if c.Quiet {
		fmt.Fprintln(out, resp.Status)
		return nil
	}

	started := "---"
	if resp.StartedAt != nil {
		started = resp.StartedAt.AsTime().String()
	}
	ended := "---"
	if resp.EndedAt != nil {
		ended = resp.EndedAt.AsTime().String()
	}

	var cpu float64
	var mem uint64
	var readBps, writeBps uint32

	if resp.Limits != nil {
		cpu = resp.Limits.CpuPercent
		mem = resp.Limits.MemoryBytes
		readBps = resp.Limits.DiskReadBps
		writeBps = resp.Limits.DiskWriteBps
	}

	data := struct {
		ID           string
		Owner        string
		Command      string
		Args         []string
		Status       string
		Started      string
		Ended        string
		ExitCode     int32
		CPULimit     float64
		MemoryLimit  uint64
		DiskReadBPS  uint32
		DiskWriteBPS uint32
	}{
		ID:           resp.JobId,
		Owner:        resp.Owner,
		Command:      resp.Command,
		Args:         resp.Args,
		Status:       resp.Status.String(),
		Started:      started,
		Ended:        ended,
		ExitCode:     resp.ExitCode,
		CPULimit:     cpu * 100,
		MemoryLimit:  mem,
		DiskReadBPS:  readBps,
		DiskWriteBPS: writeBps,
	}

	tmpl, err := template.New("status").Parse(statusTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse status template: %w", err)
	}
	if err := tmpl.Execute(out, data); err != nil {
		return fmt.Errorf("failed to execute status template: %w", err)
	}
	return nil
}
