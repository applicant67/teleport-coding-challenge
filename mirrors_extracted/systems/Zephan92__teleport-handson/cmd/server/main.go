package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/teleport-handson/api/v1"
	"github.com/teleport-handson/internal/server"
	"github.com/teleport-handson/internal/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

const port = 8080

func main() {
	slog.Info("Starting server...")

	if err := worker.SetupCgroups(); err != nil {
		slog.Error("Failed to setup cgroups", "error", err)
		os.Exit(1)
	}

	// Load Server Certificate and Key
	cert, err := tls.LoadX509KeyPair("certs/server.crt", "certs/server.key")
	if err != nil {
		slog.Error("Failed to load server key pair", "error", err)
		os.Exit(1)
	}

	// Load CA Certificate to verify clients
	caCert, err := os.ReadFile("certs/ca.crt")
	if err != nil {
		slog.Error("Failed to load CA cert", "error", err)
		os.Exit(1)
	}
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		slog.Error("Failed to append CA cert")
		os.Exit(1)
	}

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS13,
	})

	s := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterJobWorkerServer(s, server.NewWorkerServer(""))
	reflection.Register(s)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		slog.Error("Failed to listen", "error", err)
		os.Exit(1)
	}

	go func() {
		slog.Info("Job Worker Server listening", "port", port)
		if err := s.Serve(lis); err != nil {
			slog.Error("Failed to serve", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server...")
	s.GracefulStop()
	slog.Info("Server stopped")
}
