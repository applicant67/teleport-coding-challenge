package cli

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	pb "github.com/teleport-handson/api/v1"
	"google.golang.org/grpc"
)

// --- Mocks ---

type mockJobWorkerClient struct {
	pb.JobWorkerClient
	startFunc  func(ctx context.Context, in *pb.StartRequest, opts ...grpc.CallOption) (*pb.StartResponse, error)
	stopFunc   func(ctx context.Context, in *pb.StopRequest, opts ...grpc.CallOption) (*pb.StopResponse, error)
	statusFunc func(ctx context.Context, in *pb.StatusRequest, opts ...grpc.CallOption) (*pb.StatusResponse, error)
	streamFunc func(ctx context.Context, in *pb.StreamRequest, opts ...grpc.CallOption) (pb.JobWorker_StreamClient, error)
}

func (m *mockJobWorkerClient) Start(ctx context.Context, in *pb.StartRequest, opts ...grpc.CallOption) (*pb.StartResponse, error) {
	if m.startFunc != nil {
		return m.startFunc(ctx, in, opts...)
	}
	return nil, errors.New("Start not implemented")
}

func (m *mockJobWorkerClient) Stop(ctx context.Context, in *pb.StopRequest, opts ...grpc.CallOption) (*pb.StopResponse, error) {
	if m.stopFunc != nil {
		return m.stopFunc(ctx, in, opts...)
	}
	return nil, errors.New("Stop not implemented")
}

func (m *mockJobWorkerClient) Status(ctx context.Context, in *pb.StatusRequest, opts ...grpc.CallOption) (*pb.StatusResponse, error) {
	if m.statusFunc != nil {
		return m.statusFunc(ctx, in, opts...)
	}
	return nil, errors.New("Status not implemented")
}

func (m *mockJobWorkerClient) Stream(ctx context.Context, in *pb.StreamRequest, opts ...grpc.CallOption) (pb.JobWorker_StreamClient, error) {
	if m.streamFunc != nil {
		return m.streamFunc(ctx, in, opts...)
	}
	return nil, errors.New("Stream not implemented")
}

type mockStreamClient struct {
	grpc.ClientStream
	recvFunc func() (*pb.LogChunk, error)
}

func (m *mockStreamClient) Recv() (*pb.LogChunk, error) {
	if m.recvFunc != nil {
		return m.recvFunc()
	}
	return nil, io.EOF
}

// --- Tests ---

func TestStartJob(t *testing.T) {
	tests := []struct {
		name       string
		outputJSON bool
		quiet      bool
		follow     bool
		mockStart  func(ctx context.Context, in *pb.StartRequest, opts ...grpc.CallOption) (*pb.StartResponse, error)
		mockStream func(ctx context.Context, in *pb.StreamRequest, opts ...grpc.CallOption) (pb.JobWorker_StreamClient, error)
		wantOutput string
		wantErr    bool
	}{
		{
			name: "Success Text",
			mockStart: func(ctx context.Context, in *pb.StartRequest, opts ...grpc.CallOption) (*pb.StartResponse, error) {
				return &pb.StartResponse{JobId: "job-123", Status: "RUNNING"}, nil
			},
			wantOutput: "Job initiated successfully.\nID:      job-123\nStatus:  RUNNING\n",
		},
		{
			name:       "Success JSON",
			outputJSON: true,
			mockStart: func(ctx context.Context, in *pb.StartRequest, opts ...grpc.CallOption) (*pb.StartResponse, error) {
				return &pb.StartResponse{JobId: "job-123", Status: "RUNNING"}, nil
			},
			wantOutput: `"job_id": "job-123"`,
		},
		{
			name:  "Success Quiet",
			quiet: true,
			mockStart: func(ctx context.Context, in *pb.StartRequest, opts ...grpc.CallOption) (*pb.StartResponse, error) {
				return &pb.StartResponse{JobId: "job-123", Status: "RUNNING"}, nil
			},
			wantOutput: "job-123\n",
		},
		{
			name:   "Success Follow",
			follow: true,
			mockStart: func(ctx context.Context, in *pb.StartRequest, opts ...grpc.CallOption) (*pb.StartResponse, error) {
				return &pb.StartResponse{JobId: "job-123", Status: "RUNNING"}, nil
			},
			mockStream: func(ctx context.Context, in *pb.StreamRequest, opts ...grpc.CallOption) (pb.JobWorker_StreamClient, error) {
				count := 0
				return &mockStreamClient{
					recvFunc: func() (*pb.LogChunk, error) {
						if count > 0 {
							return nil, io.EOF
						}
						count++
						return &pb.LogChunk{Data: []byte("log output")}, nil
					},
				}, nil
			},
			wantOutput: "log output",
		},
		{
			name: "Start Error",
			mockStart: func(ctx context.Context, in *pb.StartRequest, opts ...grpc.CallOption) (*pb.StartResponse, error) {
				return nil, errors.New("rpc error")
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockJobWorkerClient{
				startFunc:  tc.mockStart,
				streamFunc: tc.mockStream,
			}
			cmd := &cobra.Command{}
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cfg := &Config{
				OutputJSON: tc.outputJSON,
				Quiet:      tc.quiet,
				WorkDir:    "/tmp",
				Follow:     tc.follow,
			}
			err := cfg.startJob(context.Background(), buf, client, "echo", []string{"hello"})
			if (err != nil) != tc.wantErr {
				t.Errorf("startJob() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !tc.wantErr && !strings.Contains(buf.String(), tc.wantOutput) {
				t.Errorf("Output = %q, want substring %q", buf.String(), tc.wantOutput)
			}
		})
	}
}

func TestStopJob(t *testing.T) {
	tests := []struct {
		name       string
		outputJSON bool
		quiet      bool
		mockStop   func(ctx context.Context, in *pb.StopRequest, opts ...grpc.CallOption) (*pb.StopResponse, error)
		wantOutput string
		wantErr    bool
	}{
		{
			name: "Success",
			mockStop: func(ctx context.Context, in *pb.StopRequest, opts ...grpc.CallOption) (*pb.StopResponse, error) {
				return &pb.StopResponse{Message: "stopped"}, nil
			},
			wantOutput: "Job ID 'job-123' stopped successfully",
		},
		{
			name:  "Success Quiet",
			quiet: true,
			mockStop: func(ctx context.Context, in *pb.StopRequest, opts ...grpc.CallOption) (*pb.StopResponse, error) {
				return &pb.StopResponse{Message: "stopped"}, nil
			},
			wantOutput: "job-123\n",
		},
		{
			name: "Error",
			mockStop: func(ctx context.Context, in *pb.StopRequest, opts ...grpc.CallOption) (*pb.StopResponse, error) {
				return nil, errors.New("rpc error")
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockJobWorkerClient{stopFunc: tc.mockStop}
			cmd := &cobra.Command{}
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cfg := &Config{OutputJSON: tc.outputJSON, Quiet: tc.quiet}
			err := cfg.stopJob(context.Background(), buf, client, "job-123")
			if (err != nil) != tc.wantErr {
				t.Errorf("stopJob() error = %v, wantErr %v", err, tc.wantErr)
			}
			if !tc.wantErr && !strings.Contains(buf.String(), tc.wantOutput) {
				t.Errorf("Output = %q, want substring %q", buf.String(), tc.wantOutput)
			}
		})
	}
}

func TestGetJobStatus(t *testing.T) {
	tests := []struct {
		name       string
		outputJSON bool
		quiet      bool
		mockStatus func(ctx context.Context, in *pb.StatusRequest, opts ...grpc.CallOption) (*pb.StatusResponse, error)
		wantOutput string
		wantErr    bool
	}{
		{
			name: "Success Text",
			mockStatus: func(ctx context.Context, in *pb.StatusRequest, opts ...grpc.CallOption) (*pb.StatusResponse, error) {
				return &pb.StatusResponse{
					JobId:    "job-123",
					Owner:    "alice",
					Command:  "echo",
					Args:     []string{"hello"},
					Status:   pb.StatusResponse_COMPLETED,
					ExitCode: 0,
				}, nil
			},
			wantOutput: "Job Information:\n  ID:            job-123\n  Owner:         alice\n  Command:       echo\n  Args:          [hello]\n\nLifecycle:\n  Status:        COMPLETED\n  Started:       ---\n  Ended:         ---\n  Exit Code:     0\n\nResource Constraints:\n  CPU Limit:     0.00%\n  Memory Limit:  0 bytes\n  Disk I/O:      Read: 0 bps | Write: 0 bps\n",
		},
		{
			name:       "Success JSON",
			outputJSON: true,
			mockStatus: func(ctx context.Context, in *pb.StatusRequest, opts ...grpc.CallOption) (*pb.StatusResponse, error) {
				return &pb.StatusResponse{
					JobId:    "job-123",
					Owner:    "alice",
					Command:  "echo",
					Args:     []string{"hello"},
					Status:   pb.StatusResponse_COMPLETED,
					ExitCode: 0,
				}, nil
			},
			wantOutput: `"job_id": "job-123"`,
		},
		{
			name:  "Success Quiet",
			quiet: true,
			mockStatus: func(ctx context.Context, in *pb.StatusRequest, opts ...grpc.CallOption) (*pb.StatusResponse, error) {
				return &pb.StatusResponse{
					JobId:    "job-123",
					Owner:    "alice",
					Command:  "echo",
					Args:     []string{"hello"},
					Status:   pb.StatusResponse_COMPLETED,
					ExitCode: 0,
				}, nil
			},
			wantOutput: "COMPLETED\n",
		},
		{
			name: "Error",
			mockStatus: func(ctx context.Context, in *pb.StatusRequest, opts ...grpc.CallOption) (*pb.StatusResponse, error) {
				return nil, errors.New("rpc error")
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockJobWorkerClient{statusFunc: tc.mockStatus}
			cmd := &cobra.Command{}
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cfg := &Config{OutputJSON: tc.outputJSON, Quiet: tc.quiet}
			err := cfg.getJobStatus(context.Background(), buf, client, "job-123")
			if (err != nil) != tc.wantErr {
				t.Errorf("getJobStatus() error = %v, wantErr %v", err, tc.wantErr)
			}
			if !tc.wantErr && !strings.Contains(buf.String(), tc.wantOutput) {
				t.Errorf("Output = %q, want substring %q", buf.String(), tc.wantOutput)
			}
		})
	}
}

func TestStreamLogs(t *testing.T) {
	tests := []struct {
		name       string
		mockStream func(ctx context.Context, in *pb.StreamRequest, opts ...grpc.CallOption) (pb.JobWorker_StreamClient, error)
		wantOutput string
		wantErr    bool
	}{
		{
			name: "Success",
			mockStream: func(ctx context.Context, in *pb.StreamRequest, opts ...grpc.CallOption) (pb.JobWorker_StreamClient, error) {
				count := 0
				return &mockStreamClient{
					recvFunc: func() (*pb.LogChunk, error) {
						if count == 0 {
							count++
							return &pb.LogChunk{Data: []byte("line 1\n")}, nil
						}
						if count == 1 {
							count++
							return &pb.LogChunk{Data: []byte("line 2\n")}, nil
						}
						return nil, io.EOF
					},
				}, nil
			},
			wantOutput: "line 1\nline 2\n",
		},
		{
			name: "Connect Error",
			mockStream: func(ctx context.Context, in *pb.StreamRequest, opts ...grpc.CallOption) (pb.JobWorker_StreamClient, error) {
				return nil, errors.New("connect failed")
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockJobWorkerClient{streamFunc: tc.mockStream}
			cmd := &cobra.Command{}
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cfg := &Config{}
			err := cfg.streamLogs(context.Background(), buf, client, "job-123")
			if (err != nil) != tc.wantErr {
				t.Errorf("streamLogs() error = %v, wantErr %v", err, tc.wantErr)
			}
			if !tc.wantErr && buf.String() != tc.wantOutput {
				t.Errorf("Output = %q, want %q", buf.String(), tc.wantOutput)
			}
		})
	}
}

func TestNewRootCmd(t *testing.T) {
	cmd := NewRootCmd()
	if cmd == nil {
		t.Fatal("NewRootCmd() returned nil")
	}

	if !cmd.SilenceUsage {
		t.Error("Expected SilenceUsage to be true")
	}

	// Verify flags are registered
	requiredFlags := []string{flagNameServerAddr, flagNameCAPath, flagNameCertPath, flagNameKeyPath, flagNameJSON}
	for _, flag := range requiredFlags {
		if cmd.PersistentFlags().Lookup(flag) == nil {
			t.Errorf("Missing persistent flag: %s", flag)
		}
	}

	// Verify subcommands are registered
	if len(cmd.Commands()) < 5 {
		t.Errorf("Expected at least 5 subcommands, got %d", len(cmd.Commands()))
	}
}

func TestVersionCmd(t *testing.T) {
	cmd := newVersionCommand()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	cmd.Run(cmd, []string{})
	if !strings.Contains(buf.String(), "worker-cli version:") {
		t.Errorf("Unexpected version output: %q", buf.String())
	}
}

func TestStartCommand_Wiring(t *testing.T) {
	// This test verifies that the Cobra command correctly parses flags
	// and passes them to the business logic via the mocked connection.
	mockStart := func(ctx context.Context, in *pb.StartRequest, opts ...grpc.CallOption) (*pb.StartResponse, error) {
		if in.Command != "echo" {
			t.Errorf("Expected command 'echo', got %q", in.Command)
		}
		if in.WorkingDirectory != "/custom/dir" {
			t.Errorf("Expected dir '/custom/dir', got %q", in.WorkingDirectory)
		}
		return &pb.StartResponse{JobId: "job-wiring", Status: "RUNNING"}, nil
	}

	cfg := &Config{
		Connect: func() (pb.JobWorkerClient, *grpc.ClientConn, error) {
			return &mockJobWorkerClient{startFunc: mockStart}, nil, nil
		},
	}

	cmd := newStartCommand(cfg)
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	// Simulate CLI arguments: start --dir /custom/dir -- echo hello
	cmd.SetArgs([]string{"--dir", "/custom/dir", "--", "echo", "hello"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Command execution failed: %v", err)
	}

	if !strings.Contains(buf.String(), "job-wiring") {
		t.Errorf("Output missing job ID: %s", buf.String())
	}
}
