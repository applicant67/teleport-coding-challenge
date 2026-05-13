package main

import (
	"flag"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	apiv1 "github.com/benmoss/job-worker-service/pkg/api/jobworker/v1"
	"github.com/benmoss/job-worker-service/pkg/server"
	"github.com/benmoss/job-worker-service/pkg/worker"
)

type config struct {
	addr       string
	cgroupRoot string
	logDir     string
}

func main() {
	cfg := config{}
	flag.StringVar(&cfg.addr, "addr", ":8080", "Address to listen on")
	flag.StringVar(&cfg.cgroupRoot, "cgroup-root", "/sys/fs/cgroup/jobworker", "Root directory for job cgroups")
	flag.StringVar(&cfg.logDir, "log-dir", "/var/lib/jobworker", "Parent directory for job log directories")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	if err := run(cfg); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}

func run(cfg config) error {
	w, err := worker.New(worker.Config{
		CgroupRoot: cfg.cgroupRoot,
		LogDir:     cfg.logDir,
	})
	if err != nil {
		return err
	}
	defer func() {
		if err := w.Close(); err != nil {
			slog.Error("worker close error", "error", err)
		}
	}()

	lis, err := net.Listen("tcp", cfg.addr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	apiv1.RegisterJobWorkerServiceServer(grpcServer, server.New(w))
	reflection.Register(grpcServer)

	slog.Info("starting server", "addr", cfg.addr)

	errCh := make(chan error, 1)
	go func() {
		errCh <- grpcServer.Serve(lis)
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return err
	case sig := <-sigCh:
		slog.Info("shutting down", "signal", sig)
		grpcServer.GracefulStop()
		return nil
	}
}
