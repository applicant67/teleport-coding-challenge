package worker

import (
	"os/exec"
)

type ResourceController interface {
	SetLimits(jobID string, cpuLimit int, memoryLimit int, diskIOWeight int) error
	GetLimits(jobID string) (cpuLimit int, memoryLimit int, diskIOWeight int, err error)
	CreateCgroup(jobID string) error
	StartProcessInCgroup(jobID string, cmd *exec.Cmd) error
	CleanupCgroup(jobID string) error
}
