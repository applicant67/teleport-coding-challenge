package worker

// JobLogger defines operations for logging job-related information.
type JobLogger interface {
	AddJob(jobID string)
	AddLog(jobID string, log []byte)
	GetLogs(jobID string) ([][]byte, error)
	GetOrCreateLogChannel(jobID string) chan []byte
}
