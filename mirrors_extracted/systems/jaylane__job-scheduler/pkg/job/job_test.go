package job

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	assert := assert.New(t)

	var tests = []struct {
		input    Job
		expected string
	}{
		{
			Job{
				ID: "job1",
				Process: &Process{
					Status: TERMINATED,
				},
			},
			"terminated",
		},
		{
			Job{
				ID: "job2",
				Process: &Process{
					Status: RUNNING,
				},
			},
			"running",
		},
		{
			Job{
				ID: "job3",
				Process: &Process{
					Status: STOPPED,
				},
			},
			"stopped",
		},
		{
			Job{
				ID: "job4",
				Process: &Process{
					Status: UNKNOWN,
				},
			},
			"unknown",
		}}

	for _, test := range tests {
		assert.Equalf(test.expected, test.input.Process.Status.String(), "expected String() to return %v", test.expected)
	}
}

func TestEnumIndex(t *testing.T) {
	assert := assert.New(t)

	var tests = []struct {
		input    Job
		expected int
	}{
		{
			Job{
				ID: "job1",
				Process: &Process{
					Status: TERMINATED,
				},
			},
			0,
		},
		{
			Job{
				ID: "job2",
				Process: &Process{
					Status: RUNNING,
				},
			},
			1,
		},
		{
			Job{
				ID: "job3",
				Process: &Process{
					Status: STOPPED,
				},
			},
			2,
		},
		{
			Job{
				ID: "job4",
				Process: &Process{
					Status: UNKNOWN,
				},
			},
			3,
		}}

	for _, test := range tests {
		assert.Equalf(test.expected, test.input.Process.Status.EnumIndex(), "expected String() to return %v", test.expected)
	}
}

func TestIsRunning(t *testing.T) {
	assert := assert.New(t)

	var tests = []struct {
		input    Job
		expected bool
	}{
		{
			Job{
				Process: &Process{
					ExitCode: 0,
					Status:   RUNNING,
				},
			},
			true,
		},
		{
			Job{
				Process: &Process{
					ExitCode: -1,
					Status:   TERMINATED,
				},
			},
			false,
		},
		{
			Job{
				Process: &Process{
					ExitCode: 0,
					Status:   TERMINATED,
				},
			},
			false,
		},
		{
			Job{
				Process: &Process{
					ExitCode: -1,
					Status:   RUNNING,
				},
			},
			false,
		}}

	for _, test := range tests {
		assert.Equalf(test.expected, test.input.IsRunning(), "expected IsRunning() to return %v", test.expected)
	}
}
