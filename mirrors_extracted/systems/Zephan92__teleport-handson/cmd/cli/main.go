package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/teleport-handson/internal/cli"
)

func main() {
	rootCmd := cli.NewRootCmd()

	// Handle signals (Ctrl+C) gracefully
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
