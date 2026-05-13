// Package main implements the worker-cli command-line client for the Job Worker service.
// It provides commands to start, stop, query status, and stream output from jobs.
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/alecthomas/kingpin/v2"
)

var (
	version = "dev"
	commit  = "unknown"
)

// Config holds global configuration shared across commands.
type Config struct {
	ServerAddr string
	CertFile   string
	KeyFile    string
	CAFile     string
	JSON       bool
}

// printJSON outputs v as JSON to stdout.
func printJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func main() {
	app := kingpin.New("worker-cli", "Job Worker CLI - Execute and manage jobs on a remote worker server.")
	app.UsageTemplate(kingpin.CompactUsageTemplate)
	app.Version(fmt.Sprintf("%s (commit: %s)", version, commit)).VersionFlag.Short('v')
	app.HelpFlag.Short('h')

	// Global configuration
	cfg := &Config{}
	app.Flag("json", "Output in JSON format").
		Short('j').
		BoolVar(&cfg.JSON)
	app.Flag("server", "Server address").
		Default("localhost:50051").
		Envar("WORKER_SERVER").
		HintOptions("localhost:50051", "127.0.0.1:50051").
		StringVar(&cfg.ServerAddr)
	app.Flag("cert", "Client TLS certificate path").
		Default("./certs/client1.pem").
		Envar("WORKER_CERT").
		StringVar(&cfg.CertFile)
	app.Flag("key", "Client TLS private key path").
		Default("./certs/client1-key.pem").
		Envar("WORKER_KEY").
		StringVar(&cfg.KeyFile)
	app.Flag("ca", "CA certificate path for server verification").
		Default("./certs/ca.pem").
		Envar("WORKER_CA").
		StringVar(&cfg.CAFile)

	// Register commands
	registerStartCommand(app, cfg)
	registerStopCommand(app, cfg)
	registerStatusCommand(app, cfg)
	registerLogsCommand(app, cfg)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
