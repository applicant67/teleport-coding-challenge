package main

import (
	"context"
	"fmt"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/alecthomas/kingpin/v2"
	"github.com/manil/job-worker/pkg/client"
)

// StartCommand handles the "start" subcommand.
type StartCommand struct {
	cfg      *Config
	cpuLimit float64
	memory   string
	ioWeight int32
	command  string
	args     []string
}

func registerStartCommand(app *kingpin.Application, cfg *Config) {
	c := &StartCommand{cfg: cfg}
	cmd := app.Command("start", "Start a new job").Action(c.run)
	cmd.Flag("cpu", "CPU limit as fraction (0.5 = 50% of one core, 2.0 = two cores)").
		Default("0").
		HintOptions("0.5", "1.0", "2.0").
		Float64Var(&c.cpuLimit)
	cmd.Flag("memory", "Memory limit (e.g., 512M, 1G, 1073741824)").
		Default("0").
		HintOptions("256M", "512M", "1G", "2G").
		StringVar(&c.memory)
	cmd.Flag("io-weight", "I/O weight (1-10000, 0 = default 100)").
		Default("0").
		HintOptions("50", "100", "500", "1000").
		Int32Var(&c.ioWeight)
	cmd.Arg("command", "Command to execute").
		Required().
		StringVar(&c.command)
	cmd.Arg("args", "Arguments to pass to the command").
		StringsVar(&c.args)
}

func (c *StartCommand) run(*kingpin.ParseContext) error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cli, err := client.New(c.cfg.ServerAddr, c.cfg.CertFile, c.cfg.KeyFile, c.cfg.CAFile)
	if err != nil {
		return err
	}
	defer cli.Close()

	memoryBytes, err := parseMemorySize(c.memory)
	if err != nil {
		return fmt.Errorf("invalid memory size: %w", err)
	}

	jobID, err := cli.Start(ctx, c.command, c.args,
		client.WithCPU(c.cpuLimit),
		client.WithMemory(memoryBytes),
		client.WithIOWeight(c.ioWeight),
	)
	if err != nil {
		return err
	}

	if c.cfg.JSON {
		return printJSON(map[string]string{"job_id": jobID})
	}
	fmt.Printf("Job started: %s\n", jobID)
	return nil
}

// parseMemorySize parses human-readable memory sizes like "512M", "1G", or raw bytes.
func parseMemorySize(s string) (int64, error) {
	if s == "" || s == "0" {
		return 0, nil
	}

	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}

	multiplier := int64(1)
	suffix := s[len(s)-1]

	switch suffix {
	case 'k', 'K':
		multiplier = 1024
		s = s[:len(s)-1]
	case 'm', 'M':
		multiplier = 1024 * 1024
		s = s[:len(s)-1]
	case 'g', 'G':
		multiplier = 1024 * 1024 * 1024
		s = s[:len(s)-1]
	case 't', 'T':
		multiplier = 1024 * 1024 * 1024 * 1024
		s = s[:len(s)-1]
	}

	value, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number: %s", s)
	}

	return value * multiplier, nil
}
