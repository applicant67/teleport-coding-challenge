package worker

// ResourceLimits defaults define the baseline constraints applied to jobs.
const (
	// DefaultMemoryLimitBytes is the default memory limit (100 MB).
	DefaultMemoryLimitBytes = 100 * 1024 * 1024

	// DefaultCPULimitPercent is the default CPU limit (50%).
	DefaultCPULimitPercent = 0.5

	// DefaultDiskReadBPS is the default read throughput limit in bytes per second (10 MB/s).
	DefaultDiskReadBPS = 10 * 1024 * 1024

	// DefaultDiskWriteBPS is the default write throughput limit in bytes per second (10 MB/s).
	DefaultDiskWriteBPS = 10 * 1024 * 1024
)
