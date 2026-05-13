package output

import (
	"errors"
	"testing"
)

func TestLog_Append(t *testing.T) {
	log := newTestLog(t)

	data := []byte("test data")
	n, err := log.Append(TagStdout, data)
	if err != nil {
		t.Fatalf("Append() error = %v", err)
	}

	expectedLen := int64(FrameHeaderSize + len(data))
	if int64(n) != expectedLen {
		t.Errorf("Append() = %d, want %d", n, expectedLen)
	}
	if log.committedLen != expectedLen {
		t.Errorf("CommittedLen() = %d, want %d", log.committedLen, expectedLen)
	}
}

func TestLog_AppendAfterClose(t *testing.T) {
	log := newTestLog(t)

	log.Close()

	_, err := log.Append(TagStdout, []byte("data"))
	if !errors.Is(err, ErrLogClosed) {
		t.Errorf("Append() = %v, want ErrLogClosed", err)
	}
}

func TestLog_CloseIdempotent(t *testing.T) {
	log := newTestLog(t)

	if err := log.Close(); err != nil {
		t.Fatalf("first Close() error = %v", err)
	}

	if err := log.Close(); err != nil {
		t.Errorf("second Close() error = %v", err)
	}
}

func TestLog_Subscribe(t *testing.T) {
	log := newTestLog(t)

	sub, err := log.Subscribe()
	if err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}
	defer sub.Close()
}
