//go:build !linux

package worker

import "errors"

// ErrLinuxRequired is returned on non-Linux platforms.
var ErrLinuxRequired = errors.New("job execution requires Linux with cgroups v2")

// Start is not supported on non-Linux platforms.
func (j *Job) Start(_ string) error {
	return ErrLinuxRequired
}

// Stop is not supported on non-Linux platforms.
func (j *Job) Stop() error {
	return ErrLinuxRequired
}

// Wait is not supported on non-Linux platforms.
func (j *Job) Wait() {
	// No-op
}

// Done returns a closed channel on non-Linux platforms.
func (j *Job) Done() <-chan struct{} {
	ch := make(chan struct{})
	close(ch)
	return ch
}
