package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kingpin/v2"
	"github.com/manil/job-worker/pkg/client"
)

// StopCommand handles the "stop" subcommand.
type StopCommand struct {
	cfg   *Config
	jobID string
}

func registerStopCommand(app *kingpin.Application, cfg *Config) {
	c := &StopCommand{cfg: cfg}
	cmd := app.Command("stop", "Stop a running job").Action(c.run)
	cmd.Arg("job-id", "ID of the job to stop").
		Required().
		StringVar(&c.jobID)
}

func (c *StopCommand) run(*kingpin.ParseContext) error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cli, err := client.New(c.cfg.ServerAddr, c.cfg.CertFile, c.cfg.KeyFile, c.cfg.CAFile)
	if err != nil {
		return err
	}
	defer cli.Close()

	if err := cli.Stop(ctx, c.jobID); err != nil {
		return err
	}

	if c.cfg.JSON {
		return printJSON(map[string]any{"job_id": c.jobID, "stopped": true})
	}
	fmt.Printf("Job stopped: %s\n", c.jobID)
	return nil
}
