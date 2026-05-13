package output

import (
	"bytes"
	"context"
	"errors"
	"io"
	"sync"
	"testing"
	"time"
)

func TestSubscriber_Read(t *testing.T) {
	log := newTestLog(t)

	// Append data
	data := []byte("hello world")
	if _, err := log.Append(TagStdout, data); err != nil {
		t.Fatalf("Append() error = %v", err)
	}

	// Subscribe and read
	sub, err := log.Subscribe()
	if err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}
	defer sub.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	frame, err := sub.Read(ctx)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	if frame.Tag != TagStdout {
		t.Errorf("Tag = %v, want %v", frame.Tag, TagStdout)
	}
	if !bytes.Equal(frame.Data, data) {
		t.Errorf("Data = %q, want %q", frame.Data, data)
	}
}

func TestSubscriber_LateSubscription(t *testing.T) {
	log := newTestLog(t)

	// Append BEFORE subscribing
	data1 := []byte("message 1")
	data2 := []byte("message 2")
	if _, err := log.Append(TagStdout, data1); err != nil {
		t.Fatalf("Append(data1) error = %v", err)
	}
	if _, err := log.Append(TagStderr, data2); err != nil {
		t.Fatalf("Append(data2) error = %v", err)
	}

	// Late subscriber should see all data
	sub, err := log.Subscribe()
	if err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}
	defer sub.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	frame1, err := sub.Read(ctx)
	if err != nil {
		t.Fatalf("Read(ctx) frame1 error = %v", err)
	}
	frame2, err := sub.Read(ctx)
	if err != nil {
		t.Fatalf("Read(ctx) frame2 error = %v", err)
	}

	if !bytes.Equal(frame1.Data, data1) {
		t.Errorf("frame1.Data = %q, want %q", frame1.Data, data1)
	}
	if !bytes.Equal(frame2.Data, data2) {
		t.Errorf("frame2.Data = %q, want %q", frame2.Data, data2)
	}
}

func TestSubscriber_MultipleSubscribers(t *testing.T) {
	log := newTestLog(t)

	// Create multiple subscribers
	sub1, err := log.Subscribe()
	if err != nil {
		t.Fatalf("Subscribe() sub1 error = %v", err)
	}
	defer sub1.Close()
	sub2, err := log.Subscribe()
	if err != nil {
		t.Fatalf("Subscribe() sub2 error = %v", err)
	}
	defer sub2.Close()

	// Append data
	data := []byte("shared message")
	if _, err := log.Append(TagStdout, data); err != nil {
		t.Fatalf("Append() error = %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Both should receive the data
	frame1, err1 := sub1.Read(ctx)
	frame2, err2 := sub2.Read(ctx)

	if err1 != nil || err2 != nil {
		t.Fatalf("Read errors: %v, %v", err1, err2)
	}

	if !bytes.Equal(frame1.Data, data) || !bytes.Equal(frame2.Data, data) {
		t.Error("both subscribers should receive the same data")
	}
}

func TestSubscriber_BlocksUntilData(t *testing.T) {
	log := newTestLog(t)

	sub, err := log.Subscribe()
	if err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}
	defer sub.Close()

	// Read in goroutine
	type readResult struct {
		frame *Frame
		err   error
	}
	result := make(chan readResult, 1)
	go func() {
		ctx := context.Background()
		frame, err := sub.Read(ctx)
		result <- readResult{frame: frame, err: err}
	}()

	// Ensure it's blocking
	select {
	case <-result:
		t.Fatal("Read() should block when no data")
	case <-time.After(100 * time.Millisecond):
	}

	// Now append data
	data := []byte("unblock")
	if _, err := log.Append(TagStdout, data); err != nil {
		t.Fatalf("Append() error = %v", err)
	}

	// Should unblock
	select {
	case res := <-result:
		if res.err != nil {
			t.Fatalf("Read() error = %v", res.err)
		}
		if !bytes.Equal(res.frame.Data, data) {
			t.Errorf("Data = %q, want %q", res.frame.Data, data)
		}
	case <-time.After(time.Second):
		t.Fatal("Read() did not unblock after Append()")
	}
}

func TestSubscriber_ContextCancellation(t *testing.T) {
	log := newTestLog(t)

	sub, err := log.Subscribe()
	if err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}
	defer sub.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = sub.Read(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("Read() = %v, want context.Canceled", err)
	}
}

func TestSubscriber_EOFOnClose(t *testing.T) {
	log := newTestLog(t)

	sub, err := log.Subscribe()
	if err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}
	defer sub.Close()

	// Read in goroutine
	done := make(chan error)
	go func() {
		ctx := context.Background()
		_, err := sub.Read(ctx)
		done <- err
	}()

	// Close the log
	log.Close()

	// Should return EOF
	select {
	case err := <-done:
		if err != io.EOF {
			t.Errorf("Read() = %v, want io.EOF", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Read() did not unblock after Close()")
	}
}

func TestSubscriber_ReadsAllBeforeEOF(t *testing.T) {
	log := newTestLog(t)

	// Append then close
	data := []byte("final message")
	if _, err := log.Append(TagStdout, data); err != nil {
		t.Fatalf("Append() error = %v", err)
	}
	log.Close()

	// Subscriber should get data, then EOF
	sub, err := log.Subscribe()
	if err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}
	defer sub.Close()

	ctx := context.Background()

	frame, err := sub.Read(ctx)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	if !bytes.Equal(frame.Data, data) {
		t.Errorf("Data = %q, want %q", frame.Data, data)
	}

	_, err = sub.Read(ctx)
	if err != io.EOF {
		t.Errorf("second Read() = %v, want io.EOF", err)
	}
}

func TestSubscriber_ConcurrentReadAndAppend(t *testing.T) {
	log := newTestLog(t)

	const (
		numMessages    = 100
		numSubscribers = 3
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var wg sync.WaitGroup

	// Start subscribers
	for i := 0; i < numSubscribers; i++ {
		sub, err := log.Subscribe()
		if err != nil {
			t.Fatalf("Subscribe() error = %v", err)
		}

		wg.Add(1)
		go func(sub *Subscriber, id int) {
			defer wg.Done()
			defer sub.Close()

			count := 0
			for {
				_, err := sub.Read(ctx)
				if err == io.EOF {
					break
				}
				if err != nil {
					t.Errorf("subscriber %d Read() error = %v", id, err)
					return
				}
				count++
			}

			if count != numMessages {
				t.Errorf("subscriber %d read %d messages, want %d", id, count, numMessages)
			}
		}(sub, i)
	}

	// Append messages
	for i := 0; i < numMessages; i++ {
		if _, err := log.Append(TagStdout, []byte("msg")); err != nil {
			t.Fatalf("Append() error = %v", err)
		}
	}

	// Close to signal completion
	log.Close()

	wg.Wait()
}
