//go:build linux

// Package main implements the worker-server binary, a gRPC server for job execution.
// Requires root privileges for cgroup management.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/manil/job-worker/internal/server"
	"github.com/manil/job-worker/pkg/worker"
)

var (
	version = "dev"
	commit  = "unknown"
)

// serverConfig holds the server configuration
type serverConfig struct {
	listenAddr string
	certFile   string
	keyFile    string
	caFile     string
	cgroupRoot string
}

func main() {
	app := kingpin.New("worker-server", "Job Worker Server - Execute and manage jobs with cgroups v2.")
	app.UsageTemplate(kingpin.CompactUsageTemplate)
	app.Version(fmt.Sprintf("%s (commit: %s)", version, commit)).VersionFlag.Short('v')
	app.HelpFlag.Short('h')

	cfg := &serverConfig{}

	app.Flag("listen", "Address and port to listen on").
		Default("0.0.0.0:50051").
		Envar("WORKER_LISTEN").
		HintOptions("0.0.0.0:50051", "127.0.0.1:50051", ":50051").
		StringVar(&cfg.listenAddr)
	app.Flag("cert", "Server TLS certificate path").
		Default("./certs/server.pem").
		Envar("WORKER_CERT").
		StringVar(&cfg.certFile)
	app.Flag("key", "Server TLS private key path").
		Default("./certs/server-key.pem").
		Envar("WORKER_KEY").
		StringVar(&cfg.keyFile)
	app.Flag("ca", "CA certificate path for client verification").
		Default("./certs/ca.pem").
		Envar("WORKER_CA").
		StringVar(&cfg.caFile)
	app.Flag("cgroup-root", "Base path for job cgroups").
		Default("/sys/fs/cgroup/jobworker").
		StringVar(&cfg.cgroupRoot)

	app.Action(func(*kingpin.ParseContext) error {
		return runServer(cfg)
	})

	kingpin.MustParse(app.Parse(os.Args[1:]))
}

func runServer(cfg *serverConfig) error {
	tlsConfig, err := server.NewTLSConfig(cfg.certFile, cfg.keyFile, cfg.caFile)
	if err != nil {
		return fmt.Errorf("failed to configure TLS: %w", err)
	}

	store := worker.NewJobStore()
	srv := server.NewServer(store, cfg.cgroupRoot)

	log.Printf("Starting worker-server on %s (cgroup root: %s)", cfg.listenAddr, cfg.cgroupRoot)

	return srv.ListenAndServeWithSignals(cfg.listenAddr, tlsConfig)
}
