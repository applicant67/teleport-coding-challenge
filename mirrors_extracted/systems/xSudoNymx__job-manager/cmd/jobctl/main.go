package main

import (
	"os"

	"github.com/xSudoNymx/job-manager/cmd/jobctl/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
