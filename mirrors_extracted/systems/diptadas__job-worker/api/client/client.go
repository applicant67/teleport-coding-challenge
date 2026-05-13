package client

import (
	"fmt"
	"job-worker/lib"
	"net/http"
)

// Client contains the configuration to create a http.Client to send requests to API server.
type Client struct {
	Address    string // address of the API server
	CaCert     string // root CA certificate for mTLS
	ClientCert string // client certificate for mTLS
	ClientKey  string // client key for mTLS
}

// CreateJob takes a lib.CreateJobRequest and sends a request to API server to create a new job.
// It prints the job details in StdOut.
func (c Client) CreateJob(request lib.CreateJobRequest) (string, error) {
	resp, err := c.makeAPICall("/create", http.MethodPost, request)
	if err != nil {
		fmt.Printf("Error making API request, reason: %v\n", err)
	} else {
		fmt.Println(resp)
	}
	return resp, err
}

// StopJob takes a job ID and sends a request to API server to stop the job.
// It prints the response in StdOut.
func (c Client) StopJob(jobID string) (string, error) {
	resp, err := c.makeAPICall("/stop/"+jobID, http.MethodPut, nil)
	if err != nil {
		fmt.Printf("Error making API request, reason: %v\n", err)
	} else {
		fmt.Println(resp)
	}
	return resp, err
}

// GetJobStatus takes a job ID and sends a request to API server to fetch the job status.
// It prints the response in StdOut.
func (c Client) GetJobStatus(jobID string) (string, error) {
	resp, err := c.makeAPICall("/status/"+jobID, http.MethodGet, nil)
	if err != nil {
		fmt.Printf("Error making API request, reason: %v\n", err)
	} else {
		fmt.Println(resp)
	}
	return resp, err
}
