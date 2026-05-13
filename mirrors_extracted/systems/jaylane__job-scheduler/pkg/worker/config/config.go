package config

import "os"

// Config service configuration
type Config struct {
	// LogFolder stores all job logs
	LogFolder string
	// LogChunckSize size in bytes for each log chunck read from log file
	LogChunckSize       int
	ServerAddress       string
	ServerPort          string
	ServerCertAuthority string
	ServerCert          string
	ServerPublicKey     string

	ClientCertAuthority string
	ClientCert          string
	ClientPublicKey     string
}

func NewConfig() Config {
	return Config{
		LogFolder:     os.TempDir(),
		LogChunckSize: 1024,
	}
}
