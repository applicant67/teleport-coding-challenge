package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintJSON(t *testing.T) {
	data := map[string]string{"key": "value"}
	buf := new(bytes.Buffer)

	if err := printJSON(buf, data); err != nil {
		t.Fatalf("printJSON failed: %v", err)
	}

	expected := `{
  "key": "value"
}`
	if strings.TrimSpace(buf.String()) != strings.TrimSpace(expected) {
		t.Errorf("printJSON output mismatch.\nGot:\n%s\nWant:\n%s", buf.String(), expected)
	}
}