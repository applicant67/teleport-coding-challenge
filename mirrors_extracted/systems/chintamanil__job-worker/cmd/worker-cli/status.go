package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kingpin/v2"
	"github.com/manil/job-worker/pkg/client"
)

// StatusCommand handles the "status" subcommand.
type StatusCommand struct {
	cfg   *Config
	jobID string
}

func registerStatusCommand(app *kingpin.Application, cfg *Config) {
	c := &StatusCommand{cfg: cfg}
	cmd := app.Command("status", "Get status of a job").Action(c.run)
	cmd.Arg("job-id", "ID of the job to query").
		Required().
		StringVar(&c.jobID)
}

func (c *StatusCommand) run(*kingpin.ParseContext) error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cli, err := client.New(c.cfg.ServerAddr, c.cfg.CertFile, c.cfg.KeyFile, c.cfg.CAFile)
	if err != nil {
		return err
	}
	defer cli.Close()

	status, err := cli.Status(ctx, c.jobID)
	if err != nil {
		return err
	}

	if c.cfg.JSON {
		return printJSON(map[string]any{
			"job_id":    status.ID,
			"owner":     status.Owner,
			"status":    status.Status,
			"exit_code": status.ExitCode,
			"error":     status.ErrorMessage,
			"pid":       status.PID,
			"signal":    status.Signal,
		})
	}

	fmt.Printf("Job ID:    %s\n", status.ID)
	fmt.Printf("Owner:     %s\n", status.Owner)
	fmt.Printf("Status:    %s\n", status.Status)

	switch status.Status {
	case "RUNNING":
		if status.PID > 0 {
			fmt.Printf("PID:       %d\n", status.PID)
		}
	case "COMPLETED":
		fmt.Printf("Exit Code: %d\n", status.ExitCode)
	case "FAILED":
		fmt.Printf("Exit Code: %d\n", status.ExitCode)
		if status.ErrorMessage != "" {
			fmt.Printf("Error:     %s\n", status.ErrorMessage)
		}
	case "STOPPED":
		if status.Signal > 0 {
			fmt.Printf("Signal:    %d (%s)\n", status.Signal, signalName(status.Signal))
		} else {
			fmt.Printf("Exit Code: %d\n", status.ExitCode)
		}
	}

	return nil
}

// signalName returns the name of a signal number.
func signalName(sig int32) string {
	switch sig {
	case 1:
		return "SIGHUP"
	case 2:
		return "SIGINT"
	case 3:
		return "SIGQUIT"
	case 6:
		return "SIGABRT"
	case 9:
		return "SIGKILL"
	case 15:
		return "SIGTERM"
	default:
		return fmt.Sprintf("signal %d", sig)
	}
}
