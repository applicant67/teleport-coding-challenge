package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/xSudoNymx/job-manager/internal/auth"
	"github.com/xSudoNymx/job-manager/pkg/client"
)

var rootCmd = &cobra.Command{
	Use:          "jobctl",
	Short:        "Job manager CLI",
	Long:         "A command-line interface for the job manager service.",
	SilenceUsage: true,
}

func init() {
	cobra.OnInitialize(initConfig)

	// Define flags
	rootCmd.PersistentFlags().String("addr", "localhost:9000", "Server address")
	rootCmd.PersistentFlags().String("server-name", "", "Override server name for TLS verification")
	rootCmd.PersistentFlags().String("cert", "certs/read.crt", "Client certificate")
	rootCmd.PersistentFlags().String("key", "certs/read.key", "Client private key")
	rootCmd.PersistentFlags().String("ca", "certs/ca.crt", "CA certificate")
	rootCmd.PersistentFlags().Duration("timeout", 30*time.Second, "Request timeout")

	// Bind flags to viper
	viper.BindPFlag("addr", rootCmd.PersistentFlags().Lookup("addr"))
	viper.BindPFlag("server-name", rootCmd.PersistentFlags().Lookup("server-name"))
	viper.BindPFlag("cert", rootCmd.PersistentFlags().Lookup("cert"))
	viper.BindPFlag("key", rootCmd.PersistentFlags().Lookup("key"))
	viper.BindPFlag("ca", rootCmd.PersistentFlags().Lookup("ca"))
	viper.BindPFlag("timeout", rootCmd.PersistentFlags().Lookup("timeout"))
}

func initConfig() {
	// Config file: ~/jobctl.yaml or ./jobctl.yaml
	viper.SetConfigName("jobctl")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")

	// Environment variables: JOBCTL_ADDR, JOBCTL_CERT, etc.
	viper.SetEnvPrefix("JOBCTL")
	viper.AutomaticEnv()

	// Read config file (ignore if not found)
	viper.ReadInConfig()
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

// newClient creates a new client with TLS config.
func newClient() (*client.Client, error) {
	tlsCfg, err := auth.LoadClientTLSConfig(auth.TLSConfig{
		CertFile: viper.GetString("cert"),
		KeyFile:  viper.GetString("key"),
		CAFile:   viper.GetString("ca"),
	})
	if err != nil {
		return nil, fmt.Errorf("load TLS config: %w", err)
	}

	return client.New(client.Config{
		Addr:       viper.GetString("addr"),
		TLSConfig:  tlsCfg,
		ServerName: viper.GetString("server-name"),
	})
}

// newContext creates a context with timeout.
func newContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), viper.GetDuration("timeout"))
}

// fatal prints an error and exits.
func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
	os.Exit(1)
}
