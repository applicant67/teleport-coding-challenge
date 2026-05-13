package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kingpin/v2"
	"github.com/manil/job-worker/pkg/client"
)

// LogsCommand handles the "logs" subcommand.
type LogsCommand struct {
	cfg      *Config
	follow   bool
	noFollow bool
	jobID    string
}

func registerLogsCommand(app *kingpin.Application, cfg *Config) {
	c := &LogsCommand{cfg: cfg}
	cmd := app.Command("logs", "Stream job output").Action(c.run)
	cmd.Flag("follow", "Follow output in real-time").
		Short('f').
		Default("true").
		BoolVar(&c.follow)
	cmd.Flag("no-follow", "Print current output and exit").
		BoolVar(&c.noFollow)
	cmd.Arg("job-id", "ID of the job to stream output from").
		Required().
		StringVar(&c.jobID)
}

func (c *LogsCommand) run(*kingpin.ParseContext) error {
	if c.cfg.JSON {
		return fmt.Errorf("--json is not supported for logs command (use raw output)")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cli, err := client.New(c.cfg.ServerAddr, c.cfg.CertFile, c.cfg.KeyFile, c.cfg.CAFile)
	if err != nil {
		return err
	}
	defer cli.Close()

	return cli.Logs(ctx, c.jobID, os.Stdout)
}
