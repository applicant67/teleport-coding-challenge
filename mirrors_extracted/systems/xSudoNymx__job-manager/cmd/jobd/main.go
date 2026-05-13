package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/xSudoNymx/job-manager/gen/jobpb"
	"github.com/xSudoNymx/job-manager/internal/auth"
	"github.com/xSudoNymx/job-manager/internal/server"
	"github.com/xSudoNymx/job-manager/pkg/manager"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	if err := run(); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	// Parse flags
	var (
		addr     = flag.String("addr", ":9000", "gRPC listen address")
		dataDir  = flag.String("data", "", "Data directory for job logs (default: temp)")
		certFile = flag.String("cert", "certs/server.crt", "Server certificate")
		keyFile  = flag.String("key", "certs/server.key", "Server private key")
		caFile   = flag.String("ca", "certs/ca.crt", "CA certificate")
	)

	flag.Parse()

	tlsCfg, err := auth.LoadServerTLSConfig(auth.TLSConfig{
		CertFile: *certFile,
		KeyFile:  *keyFile,
		CAFile:   *caFile,
	})
	if err != nil {
		return fmt.Errorf("load server TLS config: %w", err)
	}

	logger := slog.Default()

	// Create manager
	mgr, err := manager.New(manager.Config{
		DataDir: *dataDir,
		Logger:  logger,
	})
	if err != nil {
		return fmt.Errorf("create manager: %w", err)
	}

	defer func() {
		if err := mgr.Shutdown(); err != nil {
			slog.Warn("manager shutdown error", "error", err)
		}
	}()

	// Create gRPC server
	grpcServer := grpc.NewServer(
		grpc.Creds(credentials.NewTLS(tlsCfg)),
		grpc.UnaryInterceptor(auth.UnaryInterceptor()),
		grpc.StreamInterceptor(auth.StreamInterceptor()),
	)

	// Register service
	jobServer, err := server.NewJobGRPCServer(mgr, logger)
	if err != nil {
		return fmt.Errorf("create job grpc server: %w", err)
	}
	jobpb.RegisterJobServiceServer(grpcServer, jobServer)

	// Start listener
	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	// Handle shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	errCh := make(chan error, 1)
	go func() {
		slog.Info("server starting", "addr", *addr)
		errCh <- grpcServer.Serve(lis)
	}()

	// Wait for shutdown or error
	select {
	case err := <-errCh:
		return err
	case sig := <-shutdown:
		slog.Info("shutting down", "signal", sig)
	}

	// Graceful shutdown
	grpcServer.GracefulStop()

	slog.Info("server stopped")
	return nil
}
