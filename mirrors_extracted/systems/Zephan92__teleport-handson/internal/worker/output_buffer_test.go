package worker

import (
	"context"
	"io"
	"os"
	"sync"
	"testing"
	"time"
)

// TestOutputBuffer_WriteAndStream asserts the basic functionality of writing
// data to the buffer and having a single streaming client receive that data
// in real-time.
func TestOutputBuffer_WriteAndStream(t *testing.T) {
	buf := NewOutputBuffer()
	stream := buf.Stream(t.Context())
	defer stream.Close()

	// Write some data
	input := "Hello, World!"
	buf.Stdout().Write([]byte(input))
	buf.Close()

	// Read from stream
	entry, err := stream.Next()
	if err != nil {
		t.Fatalf("Next() failed: %v", err)
	}
	if string(entry.Data) != input {
		t.Errorf("Expected %q, got %q", input, string(entry.Data))
	}
	if entry.Source != OutputSourceStdout {
		t.Errorf("Expected source Stdout, got %v", entry.Source)
	}

	// Verify stream closes
	if _, err := stream.Next(); err != io.EOF {
		t.Errorf("Expected io.EOF, got %v", err)
	}
}

// TestOutputBuffer_History asserts that a client who starts streaming after
// data has already been written will first receive the entire history before
// getting new data.
func TestOutputBuffer_History(t *testing.T) {
	buf := NewOutputBuffer()
	inputs := []string{"Line 1", "Line 2", "Line 3"}

	// Write history
	for _, in := range inputs {
		buf.Stdout().Write([]byte(in))
	}
	buf.Close()

	// New stream should receive history
	ctx := t.Context()
	stream := buf.Stream(ctx)
	defer stream.Close()

	for _, expected := range inputs {
		entry, err := stream.Next()
		if err != nil {
			t.Fatalf("Next() failed: %v", err)
		}
		if string(entry.Data) != expected {
			t.Errorf("Expected %q, got %q", expected, string(entry.Data))
		}
	}
}

// TestOutputBuffer_Sources asserts that data written to the dedicated Stdout
// and Stderr writers is correctly tagged with the appropriate source in the
// resulting LogEntry.
func TestOutputBuffer_Sources(t *testing.T) {
	buf := NewOutputBuffer()
	stdout := buf.Stdout()
	stderr := buf.Stderr()

	io.WriteString(stdout, "out")
	io.WriteString(stderr, "err")
	buf.Close()

	stream := buf.Stream(t.Context())
	defer stream.Close()

	// Order is preserved, but we need to check sources
	count := 0
	for {
		entry, err := stream.Next()
		if err == io.EOF {
			break
		}
		count++
		txt := string(entry.Data)
		if txt == "out" && entry.Source != OutputSourceStdout {
			t.Errorf("Expected Stdout source for 'out'")
		}
		if txt == "err" && entry.Source != OutputSourceStderr {
			t.Errorf("Expected Stderr source for 'err'")
		}
	}
	if count != 2 {
		t.Errorf("Expected 2 entries, got %d", count)
	}
}

// TestOutputBuffer_MultipleStreams asserts that multiple concurrent clients can
// stream from the same buffer, and all of them receive the complete and correct
// sequence of log entries.
func TestOutputBuffer_MultipleStreams(t *testing.T) {
	buf := NewOutputBuffer()
	inputs := []string{"one", "two", "three"}
	numStreamers := 5

	// Write data in a separate goroutine
	go func() {
		for _, in := range inputs {
			time.Sleep(10 * time.Millisecond) // Introduce slight delay
			buf.Stdout().Write([]byte(in))
		}
		buf.Close()
	}()

	var wg sync.WaitGroup

	for i := 0; i < numStreamers; i++ {
		wg.Go(func() {
			stream := buf.Stream(t.Context())
			defer stream.Close()
			var received []string

			for {
				entry, err := stream.Next()
				if err == io.EOF {
					break
				}
				received = append(received, string(entry.Data))
			}

			if len(received) != len(inputs) {
				t.Errorf("Streamer %d: expected %d entries, got %d", i, len(inputs), len(received))
				return
			}
			for j, val := range received {
				if val != inputs[j] {
					t.Errorf("Streamer %d: mismatch at index %d. Got %q, want %q", i, j, val, inputs[j])
				}
			}
		})
	}

	wg.Wait()
}

// TestOutputBuffer_Errors covers failure scenarios for the OutputBuffer.
func TestOutputBuffer_Errors(t *testing.T) {
	t.Run("WriteAfterClose", func(t *testing.T) {
		buf := NewOutputBuffer()
		buf.Close()
		if _, err := buf.Stdout().Write([]byte("fail")); err != os.ErrClosed {
			t.Error("Expected error writing to closed buffer")
		}
	})

	t.Run("StreamContextCancel", func(t *testing.T) {
		buf := NewOutputBuffer()
		ctx, cancel := context.WithCancel(t.Context())
		stream := buf.Stream(ctx)
		defer stream.Close()

		// Cancel immediately
		cancel()

		if _, err := stream.Next(); err != context.Canceled {
			t.Errorf("Expected context.Canceled, got %v", err)
		}
	})
}
