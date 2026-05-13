package api

import (
	"encoding/json"
	"job-worker/api/client"
	"job-worker/api/server"
	"job-worker/lib"
	"strings"
	"testing"
	"time"
)

func TestServerClientConnection(t *testing.T) {
	apiServer := server.Server{
		Port:       8443,
		CaCert:     "../ssl/ca.crt",
		ServerCert: "../ssl/server.crt",
		ServerKey:  "../ssl/server.key",
	}

	apiClient := client.Client{
		Address:    "https://localhost:8443",
		CaCert:     "../ssl/ca.crt",
		ClientCert: "../ssl/client-alice.crt",
		ClientKey:  "../ssl/client-alice.key",
	}

	go apiServer.InitializeAndRun()

	// wait to start the server
	time.Sleep(time.Second)

	// we don't need to create a job to check connection
	// just check the status for a invalid job id
	expectedResp := `{"error":"job 1: not found"}`
	if resp, err := apiClient.GetJobStatus("1"); err != nil {
		t.Fatalf("Found error: %v", err)
	} else if resp != expectedResp {
		t.Fatalf("Response not mathced, expected: %v, found: %v", expectedResp, resp)
	}
}

func TestInvalidClientCert(t *testing.T) {
	apiServer := server.Server{
		Port:       6443,
		CaCert:     "../ssl/ca.crt",
		ServerCert: "../ssl/server.crt",
		ServerKey:  "../ssl/server.key",
	}

	apiClient := client.Client{
		Address:    "https://localhost:6443",
		CaCert:     "../ssl/ca.crt",
		ClientCert: "../ssl/client-invalid.crt",
		ClientKey:  "../ssl/client-invalid.key",
	}

	go apiServer.InitializeAndRun()

	// wait to start the server
	time.Sleep(time.Second)

	// we don't need to create a job to check connection
	// just check the status for a invalid job id
	expectedError := "tls: bad certificate"
	if _, err := apiClient.GetJobStatus("1"); err == nil || !strings.Contains(err.Error(), expectedError) {
		t.Fatalf("Error not mathced, expected: %v, found: %v", expectedError, err)
	}
}

func TestInvalidUsername(t *testing.T) {
	apiServer := server.Server{
		Port:       4443,
		CaCert:     "../ssl/ca.crt",
		ServerCert: "../ssl/server.crt",
		ServerKey:  "../ssl/server.key",
	}

	apiClient := client.Client{
		Address:    "https://localhost:4443",
		CaCert:     "../ssl/ca.crt",
		ClientCert: "../ssl/client-unknown.crt",
		ClientKey:  "../ssl/client-unknown.key",
	}

	go apiServer.InitializeAndRun()

	// wait to start the server
	time.Sleep(time.Second)

	// we don't need to create a job to check user authentication
	// just check the status for a invalid job id
	expectedResp := `{"error":"user unknown not found"}`
	if resp, err := apiClient.GetJobStatus("1"); err != nil {
		t.Fatalf("Found connection error: %v", err)
	} else if !strings.Contains(resp, expectedResp) {
		t.Fatalf("Response not mathced, expected: %v, found: %v", expectedResp, resp)
	}
}

func TestInvalidPermission(t *testing.T) {
	apiServer := server.Server{
		Port:       9443,
		CaCert:     "../ssl/ca.crt",
		ServerCert: "../ssl/server.crt",
		ServerKey:  "../ssl/server.key",
	}

	apiClient := client.Client{
		Address:    "https://localhost:9443",
		CaCert:     "../ssl/ca.crt",
		ClientCert: "../ssl/client-bob.crt",
		ClientKey:  "../ssl/client-bob.key",
	}

	go apiServer.InitializeAndRun()

	// wait to start the server
	time.Sleep(time.Second)

	request := lib.CreateJobRequest{
		Command: "pwd",
	}

	// user bob is not allowed to create a job
	expectedResp := `{"error":"user bob do not have permission READ_WRITE"}`
	if resp, err := apiClient.CreateJob(request); err != nil {
		t.Fatalf("Found connection error: %v", err)
	} else if !strings.Contains(resp, expectedResp) {
		t.Fatalf("Response not mathced, expected: %v, found: %v", expectedResp, resp)
	}
}

func TestOutput(t *testing.T) {
	apiServer := server.Server{
		Port:       3443,
		CaCert:     "../ssl/ca.crt",
		ServerCert: "../ssl/server.crt",
		ServerKey:  "../ssl/server.key",
	}

	apiClient := client.Client{
		Address:    "https://localhost:3443",
		CaCert:     "../ssl/ca.crt",
		ClientCert: "../ssl/client-alice.crt",
		ClientKey:  "../ssl/client-alice.key",
	}

	go apiServer.InitializeAndRun()

	// wait to start the server
	time.Sleep(time.Second)

	request := lib.CreateJobRequest{
		Command: "echo",
		Args:    []string{"hello", "world"},
	}

	resp, err := apiClient.CreateJob(request)
	if err != nil {
		t.Fatalf("Found error: %v", err)
	} else if resp == "" {
		t.Fatal("Empty response")
	}

	var job lib.Job
	if err = json.Unmarshal([]byte(resp), &job); err != nil {
		t.Fatal("Failed to unmarshal job object")
	}

	resp, err = apiClient.GetJobStatus(job.ID)
	if err != nil {
		t.Fatalf("Found error: %v", err)
	} else if resp == "" {
		t.Fatal("Empty response")
	}

	expectedOutput := "hello world\n"

	if err = json.Unmarshal([]byte(resp), &job); err != nil {
		t.Fatal("Failed to unmarshal job object")
	} else if job.Output != expectedOutput {
		t.Fatalf("Output not mathced, expected: %v, found: %v", expectedOutput, job.Output)
	}
}
