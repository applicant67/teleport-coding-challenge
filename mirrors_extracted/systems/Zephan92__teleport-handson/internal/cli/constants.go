package cli

const (
	// Defaults for CLI flags.
	defaultServerAddr = "localhost:8080"
	defaultCAPath     = "certs/ca.crt"
	defaultCertPath   = "certs/alice.crt"
	defaultKeyPath    = "certs/alice.key"
	defaultWorkDir    = "/tmp"

	// Environment Variables used for configuration overrides.
	envVarServerAddr = "WORKER_ADDR"
	envVarCAPath     = "WORKER_CA"
	envVarCertPath   = "WORKER_CERT"
	envVarKeyPath    = "WORKER_KEY"

	// Flags used in the CLI commands.
	flagNameServerAddr = "addr"
	flagNameCAPath     = "ca"
	flagNameCertPath   = "cert"
	flagNameKeyPath    = "key"
	flagNameJSON       = "json"
	flagNameWorkDir    = "dir"
	flagNameFollow     = "follow"
	flagNameQuiet      = "quiet"

	// Usage Messages for flags.
	helpServerAddr = "Server address"
	helpCAPath     = "Path to CA certificate"
	helpCertPath   = "Path to client certificate"
	helpKeyPath    = "Path to client key"
	helpJSON       = "Output response in JSON format"
	helpWorkDir    = "Working directory for the job"
	helpFollow     = "Stream logs after starting the job"
	helpQuiet      = "Suppress informational output"
)
