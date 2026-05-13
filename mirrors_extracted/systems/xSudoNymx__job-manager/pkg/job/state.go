package job

import "fmt"

// State represents the current state of a job in its lifecycle.
type State int

const (
	Unspecified State = iota
	Running
	Stopping
	Stopped
	Exited
)

func (s State) String() string {
	switch s {
	case Unspecified:
		return "unspecified"
	case Running:
		return "running"
	case Stopping:
		return "stopping"
	case Stopped:
		return "stopped"
	case Exited:
		return "exited"
	default:
		return fmt.Sprintf("unknown(%d)", s)
	}
}
