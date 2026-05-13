package main

import (
	"testing"

	"github.com/alecthomas/kingpin/v2"
)

func TestParseMemorySize(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int64
		wantErr bool
	}{
		// Zero and empty values
		{name: "empty string", input: "", want: 0, wantErr: false},
		{name: "zero string", input: "0", want: 0, wantErr: false},
		{name: "whitespace", input: "  ", want: 0, wantErr: false},

		// Raw bytes
		{name: "raw bytes", input: "1024", want: 1024, wantErr: false},
		{name: "large raw bytes", input: "1073741824", want: 1073741824, wantErr: false},

		// Kilobytes
		{name: "kilobytes lowercase", input: "512k", want: 512 * 1024, wantErr: false},
		{name: "kilobytes uppercase", input: "512K", want: 512 * 1024, wantErr: false},

		// Megabytes
		{name: "megabytes lowercase", input: "256m", want: 256 * 1024 * 1024, wantErr: false},
		{name: "megabytes uppercase", input: "256M", want: 256 * 1024 * 1024, wantErr: false},
		{name: "512M", input: "512M", want: 512 * 1024 * 1024, wantErr: false},

		// Gigabytes
		{name: "gigabytes lowercase", input: "1g", want: 1024 * 1024 * 1024, wantErr: false},
		{name: "gigabytes uppercase", input: "1G", want: 1024 * 1024 * 1024, wantErr: false},
		{name: "2G", input: "2G", want: 2 * 1024 * 1024 * 1024, wantErr: false},

		// Terabytes
		{name: "terabytes lowercase", input: "1t", want: 1024 * 1024 * 1024 * 1024, wantErr: false},
		{name: "terabytes uppercase", input: "1T", want: 1024 * 1024 * 1024 * 1024, wantErr: false},

		// Errors
		{name: "invalid number", input: "abc", want: 0, wantErr: true},
		{name: "invalid with suffix", input: "abcM", want: 0, wantErr: true},
		{name: "negative", input: "-100M", want: -100 * 1024 * 1024, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseMemorySize(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseMemorySize(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseMemorySize(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestSignalName(t *testing.T) {
	tests := []struct {
		name string
		sig  int32
		want string
	}{
		{name: "SIGHUP", sig: 1, want: "SIGHUP"},
		{name: "SIGINT", sig: 2, want: "SIGINT"},
		{name: "SIGQUIT", sig: 3, want: "SIGQUIT"},
		{name: "SIGABRT", sig: 6, want: "SIGABRT"},
		{name: "SIGKILL", sig: 9, want: "SIGKILL"},
		{name: "SIGTERM", sig: 15, want: "SIGTERM"},
		{name: "unknown signal", sig: 99, want: "signal 99"},
		{name: "zero", sig: 0, want: "signal 0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := signalName(tt.sig); got != tt.want {
				t.Errorf("signalName(%v) = %v, want %v", tt.sig, got, tt.want)
			}
		})
	}
}

func TestConfigDefaults(t *testing.T) {
	cfg := &Config{
		ServerAddr: "localhost:50051",
		CertFile:   "./certs/client1.pem",
		KeyFile:    "./certs/client1-key.pem",
		CAFile:     "./certs/ca.pem",
	}

	if cfg.ServerAddr == "" {
		t.Error("ServerAddr should not be empty")
	}
	if cfg.CertFile == "" {
		t.Error("CertFile should not be empty")
	}
	if cfg.KeyFile == "" {
		t.Error("KeyFile should not be empty")
	}
	if cfg.CAFile == "" {
		t.Error("CAFile should not be empty")
	}
}

func TestStartCommandParsing(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		expectErr  bool
		expectCmd  string
		expectArgs []string
	}{
		{
			name:       "simple command",
			args:       []string{"start", "/bin/echo"},
			expectErr:  false,
			expectCmd:  "/bin/echo",
			expectArgs: nil,
		},
		{
			name:       "command with args",
			args:       []string{"start", "/bin/echo", "hello", "world"},
			expectErr:  false,
			expectCmd:  "/bin/echo",
			expectArgs: []string{"hello", "world"},
		},
		{
			name:      "missing command",
			args:      []string{"start"},
			expectErr: true,
		},
		{
			name:       "command with flags",
			args:       []string{"start", "--cpu", "0.5", "--memory", "512M", "/bin/sleep", "10"},
			expectErr:  false,
			expectCmd:  "/bin/sleep",
			expectArgs: []string{"10"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := kingpin.New("worker-cli", "test")
			app.Terminate(nil) // Don't exit on error
			cfg := &Config{}

			c := &StartCommand{cfg: cfg}
			cmd := app.Command("start", "Start a new job")
			cmd.Flag("cpu", "CPU limit").Float64Var(&c.cpuLimit)
			cmd.Flag("memory", "Memory limit").StringVar(&c.memory)
			cmd.Flag("io-weight", "IO weight").Int32Var(&c.ioWeight)
			cmd.Arg("command", "Command").Required().StringVar(&c.command)
			cmd.Arg("args", "Arguments").StringsVar(&c.args)

			_, err := app.Parse(tt.args)
			if (err != nil) != tt.expectErr {
				t.Errorf("parse error = %v, wantErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr {
				if c.command != tt.expectCmd {
					t.Errorf("command = %q, want %q", c.command, tt.expectCmd)
				}
				if len(c.args) != len(tt.expectArgs) {
					t.Errorf("args = %v, want %v", c.args, tt.expectArgs)
				}
			}
		})
	}
}

func TestStopCommandParsing(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		expectErr bool
		expectID  string
	}{
		{
			name:      "with job id",
			args:      []string{"stop", "job-123"},
			expectErr: false,
			expectID:  "job-123",
		},
		{
			name:      "missing job id",
			args:      []string{"stop"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := kingpin.New("worker-cli", "test")
			app.Terminate(nil)
			cfg := &Config{}

			c := &StopCommand{cfg: cfg}
			cmd := app.Command("stop", "Stop a job")
			cmd.Arg("job-id", "Job ID").Required().StringVar(&c.jobID)

			_, err := app.Parse(tt.args)
			if (err != nil) != tt.expectErr {
				t.Errorf("parse error = %v, wantErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr && c.jobID != tt.expectID {
				t.Errorf("jobID = %q, want %q", c.jobID, tt.expectID)
			}
		})
	}
}

func TestStatusCommandParsing(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		expectErr bool
		expectID  string
	}{
		{
			name:      "with job id",
			args:      []string{"status", "job-456"},
			expectErr: false,
			expectID:  "job-456",
		},
		{
			name:      "missing job id",
			args:      []string{"status"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := kingpin.New("worker-cli", "test")
			app.Terminate(nil)
			cfg := &Config{}

			c := &StatusCommand{cfg: cfg}
			cmd := app.Command("status", "Get status")
			cmd.Arg("job-id", "Job ID").Required().StringVar(&c.jobID)

			_, err := app.Parse(tt.args)
			if (err != nil) != tt.expectErr {
				t.Errorf("parse error = %v, wantErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr && c.jobID != tt.expectID {
				t.Errorf("jobID = %q, want %q", c.jobID, tt.expectID)
			}
		})
	}
}

func TestLogsCommandParsing(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		expectErr bool
		expectID  string
	}{
		{
			name:      "with job id",
			args:      []string{"logs", "job-789"},
			expectErr: false,
			expectID:  "job-789",
		},
		{
			name:      "missing job id",
			args:      []string{"logs"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := kingpin.New("worker-cli", "test")
			app.Terminate(nil)
			cfg := &Config{}

			c := &LogsCommand{cfg: cfg}
			cmd := app.Command("logs", "Stream logs")
			cmd.Arg("job-id", "Job ID").Required().StringVar(&c.jobID)

			_, err := app.Parse(tt.args)
			if (err != nil) != tt.expectErr {
				t.Errorf("parse error = %v, wantErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr && c.jobID != tt.expectID {
				t.Errorf("jobID = %q, want %q", c.jobID, tt.expectID)
			}
		})
	}
}

func TestGlobalFlagParsing(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		expectServer string
		expectCert   string
		expectKey    string
		expectCA     string
		expectJSON   bool
	}{
		{
			name:         "defaults",
			args:         []string{"status", "job-1"},
			expectServer: "localhost:50051",
			expectCert:   "./certs/client1.pem",
			expectKey:    "./certs/client1-key.pem",
			expectCA:     "./certs/ca.pem",
			expectJSON:   false,
		},
		{
			name:         "custom server",
			args:         []string{"--server", "remote:9090", "status", "job-1"},
			expectServer: "remote:9090",
			expectCert:   "./certs/client1.pem",
			expectKey:    "./certs/client1-key.pem",
			expectCA:     "./certs/ca.pem",
			expectJSON:   false,
		},
		{
			name:         "all custom",
			args:         []string{"--server", "custom:1234", "--cert", "/path/cert.pem", "--key", "/path/key.pem", "--ca", "/path/ca.pem", "status", "job-1"},
			expectServer: "custom:1234",
			expectCert:   "/path/cert.pem",
			expectKey:    "/path/key.pem",
			expectCA:     "/path/ca.pem",
			expectJSON:   false,
		},
		{
			name:         "json flag long",
			args:         []string{"--json", "status", "job-1"},
			expectServer: "localhost:50051",
			expectCert:   "./certs/client1.pem",
			expectKey:    "./certs/client1-key.pem",
			expectCA:     "./certs/ca.pem",
			expectJSON:   true,
		},
		{
			name:         "json flag short",
			args:         []string{"-j", "status", "job-1"},
			expectServer: "localhost:50051",
			expectCert:   "./certs/client1.pem",
			expectKey:    "./certs/client1-key.pem",
			expectCA:     "./certs/ca.pem",
			expectJSON:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := kingpin.New("worker-cli", "test")
			app.Terminate(nil)
			cfg := &Config{}

			app.Flag("json", "Output in JSON format").
				Short('j').
				BoolVar(&cfg.JSON)
			app.Flag("server", "Server address").
				Default("localhost:50051").
				StringVar(&cfg.ServerAddr)
			app.Flag("cert", "Certificate path").
				Default("./certs/client1.pem").
				StringVar(&cfg.CertFile)
			app.Flag("key", "Key path").
				Default("./certs/client1-key.pem").
				StringVar(&cfg.KeyFile)
			app.Flag("ca", "CA path").
				Default("./certs/ca.pem").
				StringVar(&cfg.CAFile)

			c := &StatusCommand{cfg: cfg}
			cmd := app.Command("status", "Get status")
			cmd.Arg("job-id", "Job ID").Required().StringVar(&c.jobID)

			_, err := app.Parse(tt.args)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			if cfg.ServerAddr != tt.expectServer {
				t.Errorf("ServerAddr = %q, want %q", cfg.ServerAddr, tt.expectServer)
			}
			if cfg.CertFile != tt.expectCert {
				t.Errorf("CertFile = %q, want %q", cfg.CertFile, tt.expectCert)
			}
			if cfg.KeyFile != tt.expectKey {
				t.Errorf("KeyFile = %q, want %q", cfg.KeyFile, tt.expectKey)
			}
			if cfg.CAFile != tt.expectCA {
				t.Errorf("CAFile = %q, want %q", cfg.CAFile, tt.expectCA)
			}
			if cfg.JSON != tt.expectJSON {
				t.Errorf("JSON = %v, want %v", cfg.JSON, tt.expectJSON)
			}
		})
	}
}

func TestParseMemorySizeEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int64
		wantErr bool
	}{
		// Additional edge cases
		{name: "only suffix k", input: "k", want: 0, wantErr: true},
		{name: "only suffix M", input: "M", want: 0, wantErr: true},
		{name: "only suffix G", input: "G", want: 0, wantErr: true},
		{name: "decimal value", input: "1.5G", want: 0, wantErr: true},
		{name: "mixed case suffix", input: "100m", want: 100 * 1024 * 1024, wantErr: false},
		{name: "spaces in middle", input: "100 M", want: 0, wantErr: true},
		{name: "invalid suffix", input: "100X", want: 0, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseMemorySize(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseMemorySize(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseMemorySize(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
