package lib

import (
	"strings"
	"testing"
	"time"
)

func TestCreateJobError(t *testing.T) {
	jobWorker := NewJobWorker()

	request := CreateJobRequest{
		Command: "unknown",
	}

	expectedError := "executable file not found"

	if _, err := jobWorker.CreateJob(request); err == nil || !strings.Contains(err.Error(), expectedError) {
		t.Fatalf("expected error: %v, found: %v", expectedError, err)
	}
}

func TestJobNotFoundError(t *testing.T) {
	jobWorker := NewJobWorker()

	expectedError := "not found"

	if _, err := jobWorker.GetJobStatus("0"); err == nil || !strings.Contains(err.Error(), expectedError) {
		t.Fatalf("expected error: %v, found: %v", expectedError, err)
	}
}

func TestJobOutput(t *testing.T) {
	jobWorker := NewJobWorker()

	request := CreateJobRequest{
		Command: "echo",
		Args:    []string{"hello world"},
	}

	job, err := jobWorker.CreateJob(request)
	if err != nil {
		t.Fatalf("found error: %v", err)
	}

	// wait for job to finish and copy output
	time.Sleep(time.Second)

	expectedStatus := JobExited
	expectedOutput := "hello world\n"

	job, err = jobWorker.GetJobStatus(job.ID)
	if err != nil {
		t.Fatalf("found error: %v", err)
	} else if job.Status != expectedStatus {
		t.Fatalf("expected status: %v, found: %v", expectedStatus, job.Status)
	} else if job.Output != expectedOutput {
		t.Fatalf("expected output: %v, found: %v", expectedOutput, job.Output)
	}
}

func TestForceStop(t *testing.T) {
	jobWorker := NewJobWorker()

	request := CreateJobRequest{
		Command: "sleep",
		Args:    []string{"300"},
	}

	job, err := jobWorker.CreateJob(request)
	if err != nil {
		t.Fatalf("found error: %v", err)
	}

	if err = jobWorker.StopJob(job.ID); err != nil { // failed to terminate job
		t.Fatalf("found error: %v", err)
	}

	// job terminated, check status
	expectedStatus := JobExited
	expectedError := "signal: killed"

	job, err = jobWorker.GetJobStatus(job.ID)
	if err != nil {
		t.Fatalf("found error: %v", err)
	} else if job.Status != expectedStatus {
		t.Fatalf("expected status: %v, found: %v", expectedStatus, job.Status)
	} else if job.Error == nil || !strings.Contains(job.Error.Error(), expectedError) {
		t.Fatalf("expected job error: %v, found: %v", expectedError, job.Error)
	}
}

func TestRunningJobOutput(t *testing.T) {
	jobWorker := NewJobWorker()

	request := CreateJobRequest{
		Command: "bash",
		Args:    []string{"-c", "for i in {0..100}; do echo $i; sleep 0.5; done;"},
	}

	job, err := jobWorker.CreateJob(request)
	if err != nil {
		t.Fatalf("found error: %v", err)
	}

	// wait for generating some output
	time.Sleep(time.Second)

	expectedStatus := JobRunning

	job, err = jobWorker.GetJobStatus(job.ID)
	if err != nil {
		t.Fatalf("found error: %v", err)
	} else if job.Status != expectedStatus {
		t.Fatalf("expected status: %v, found: %v", expectedStatus, job.Status)
	} else if job.Output == "" {
		t.Fatalf("found empty output")
	}

	oldOutput := job.Output
	t.Log(oldOutput)

	// get output again to ensure it contains the whole output, not just the unread part
	// wait for generating some new output
	time.Sleep(time.Second)

	job, err = jobWorker.GetJobStatus(job.ID)
	if err != nil {
		t.Fatalf("found error: %v", err)
	} else if job.Status != expectedStatus {
		t.Fatalf("expected status: %v, found: %v", expectedStatus, job.Status)
	} else if job.Output == "" || !strings.HasPrefix(job.Output, oldOutput) { // oldOutput should be prefix of current output
		t.Fatalf("found empty output")
	}

	t.Log(job.Output)
}
