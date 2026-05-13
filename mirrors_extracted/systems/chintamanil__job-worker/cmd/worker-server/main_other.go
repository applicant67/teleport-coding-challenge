//go:build !linux

// Package main provides a stub for non-Linux systems.
// The worker-server requires Linux with cgroups v2 for job execution.
package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin/v2"
)

var (
	version = "dev"
	commit  = "unknown"
)

func main() {
	app := kingpin.New("worker-server", "Job Worker Server - Execute and manage jobs with cgroups v2.")
	app.UsageTemplate(kingpin.CompactUsageTemplate)
	app.Version(fmt.Sprintf("%s (commit: %s)", version, commit)).VersionFlag.Short('v')
	app.HelpFlag.Short('h')

	app.Flag("listen", "Address and port to listen on").
		Default("0.0.0.0:50051").
		Envar("WORKER_LISTEN").
		HintOptions("0.0.0.0:50051", "127.0.0.1:50051", ":50051").
		String()
	app.Flag("cert", "Server TLS certificate path").
		Default("./certs/server.pem").
		Envar("WORKER_CERT").
		String()
	app.Flag("key", "Server TLS private key path").
		Default("./certs/server-key.pem").
		Envar("WORKER_KEY").
		String()
	app.Flag("ca", "CA certificate path for client verification").
		Default("./certs/ca.pem").
		Envar("WORKER_CA").
		String()
	app.Flag("cgroup-root", "Base path for job cgroups").
		Default("/sys/fs/cgroup/jobworker").
		String()

	app.Action(func(*kingpin.ParseContext) error {
		fmt.Fprintln(os.Stderr, "error: worker-server requires Linux with cgroups v2")
		fmt.Fprintln(os.Stderr, "       This binary cannot run on this platform.")
		os.Exit(1)
		return nil
	})

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
