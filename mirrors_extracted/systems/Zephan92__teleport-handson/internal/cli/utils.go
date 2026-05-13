// Package cli implements the command-line interface for the Job Worker Service.
package cli

import (
	"encoding/json"
	"fmt"
	"io"
)

// printJSON encodes the provided value as indented JSON and writes it to the writer.
// It is used for machine-readable output mode.
func printJSON(w io.Writer, v interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}
	return nil
}
