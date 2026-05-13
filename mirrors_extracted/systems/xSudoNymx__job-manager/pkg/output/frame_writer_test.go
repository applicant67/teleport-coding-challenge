package output

import (
	"bytes"
	"context"
	"testing"
	"time"
)

func TestFrameWriter(t *testing.T) {
	log := newTestLog(t)

	stdout := NewFrameWriter(log, TagStdout)
	stderr := NewFrameWriter(log, TagStderr)

	if _, err := stdout.Write([]byte("out")); err != nil {
		t.Fatalf("stdout.Write() error = %v", err)
	}
	if _, err := stderr.Write([]byte("err")); err != nil {
		t.Fatalf("stderr.Write() error = %v", err)
	}

	sub, err := log.Subscribe()
	if err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}
	defer sub.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	f1, err := sub.Read(ctx)
	if err != nil {
		t.Fatalf("Read() f1 error = %v", err)
	}
	f2, err := sub.Read(ctx)
	if err != nil {
		t.Fatalf("Read() f2 error = %v", err)
	}

	if f1.Tag != TagStdout || string(f1.Data) != "out" {
		t.Errorf("frame1 = %v/%q, want stdout/out", f1.Tag, f1.Data)
	}
	if f2.Tag != TagStderr || string(f2.Data) != "err" {
		t.Errorf("frame2 = %v/%q, want stderr/err", f2.Tag, f2.Data)
	}
}

func TestFrameWriter_LargeWrite(t *testing.T) {
	log := newTestLog(t)

	writer := NewFrameWriter(log, TagStdout)

	// Write more than MaxPayloadSize
	data := make([]byte, MaxPayloadSize*2+100)
	for i := range data {
		data[i] = byte(i % 256)
	}

	n, err := writer.Write(data)
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if n != len(data) {
		t.Errorf("Write() = %d, want %d", n, len(data))
	}

	// Read all chunks
	sub, err := log.Subscribe()
	if err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}
	defer sub.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var received []byte
	for len(received) < len(data) {
		frame, err := sub.Read(ctx)
		if err != nil {
			t.Fatalf("Read() error = %v", err)
		}
		received = append(received, frame.Data...)
	}

	if !bytes.Equal(received, data) {
		t.Error("received data does not match written data")
	}
}

func TestFrameWriter_EmptyWrite(t *testing.T) {
	log := newTestLog(t)

	writer := NewFrameWriter(log, TagStdout)

	n, err := writer.Write([]byte{})
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if n != 0 {
		t.Errorf("Write() = %d, want 0", n)
	}
}
