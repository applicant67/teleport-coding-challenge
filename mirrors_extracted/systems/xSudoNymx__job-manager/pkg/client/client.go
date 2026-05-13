package client

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/xSudoNymx/job-manager/gen/jobpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Client wraps the gRPC job service client.
type Client struct {
	conn *grpc.ClientConn
	svc  jobpb.JobServiceClient
}

// Config holds client configuration.
type Config struct {
	Addr       string
	TLSConfig  *tls.Config
	ServerName string
}

func New(cfg Config) (*Client, error) {
	if cfg.TLSConfig == nil {
		return nil, fmt.Errorf("tls config is required")
	}

	tlsCfg := cfg.TLSConfig.Clone()

	if cfg.ServerName != "" {
		tlsCfg.ServerName = cfg.ServerName
	}

	transportCredentials := credentials.NewTLS(tlsCfg)
	conn, err := grpc.NewClient(cfg.Addr, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}

	return &Client{
		conn: conn,
		svc:  jobpb.NewJobServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

// Start starts a new job.
func (c *Client) Start(ctx context.Context, exe string, args []string) (*jobpb.StartResponse, error) {
	return c.svc.Start(ctx, &jobpb.StartRequest{
		Spec: &jobpb.JobSpec{
			Exe:  exe,
			Args: args,
		},
	})
}

// Stop stops a running job.
func (c *Client) Stop(ctx context.Context, jobID string) (*jobpb.StopResponse, error) {
	return c.svc.Stop(ctx, &jobpb.StopRequest{
		JobId: jobID,
	})
}

// Status returns the current status of a job.
func (c *Client) Status(ctx context.Context, jobID string) (*jobpb.StatusResponse, error) {
	return c.svc.Status(ctx, &jobpb.StatusRequest{
		JobId: jobID,
	})
}

// StreamOutput returns a stream of job output.
//
// TODO: Consider returning io.ReadCloser instead of the raw gRPC stream
// to decouple consumers from the protobuf layer
func (c *Client) StreamOutput(ctx context.Context, jobID string) (jobpb.JobService_StreamOutputClient, error) {
	return c.svc.StreamOutput(ctx, &jobpb.StreamOutputRequest{
		JobId: jobID,
	})
}
