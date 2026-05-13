package worker

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestMemoryBuffer_WriteRead(t *testing.T) {
	buf := NewMemoryBuffer()
	data := []byte("hello world")

	n, err := buf.Write(data)
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if n != len(data) {
		t.Fatalf("Write() n = %d, want %d", n, len(data))
	}

	buf.Close()

	reader := buf.NewReader(context.Background())
	defer reader.Close()

	got := make([]byte, len(data))
	n, err = reader.Read(got)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	if n != len(data) {
		t.Fatalf("Read() n = %d, want %d", n, len(data))
	}
	if !bytes.Equal(got, data) {
		t.Fatalf("Read() got = %v, want %v", got, data)
	}
}

func TestMemoryBuffer_ReadFromStart(t *testing.T) {
	buf := NewMemoryBuffer()
	data1 := []byte("first ")
	data2 := []byte("second")

	if _, err := buf.Write(data1); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if _, err := buf.Write(data2); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	buf.Close()

	// Reader created after writes should get all data from start
	reader := buf.NewReader(context.Background())
	defer reader.Close()

	got, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}
	want := append(data1, data2...)
	if !bytes.Equal(got, want) {
		t.Fatalf("ReadAll() got = %q, want %q", got, want)
	}
}

func TestMemoryBuffer_ConcurrentReaders(t *testing.T) {
	buf := NewMemoryBuffer()
	data := []byte("shared data")

	if _, err := buf.Write(data); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	buf.Close()

	// Multiple readers should each get all data
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			reader := buf.NewReader(context.Background())
			defer reader.Close()

			got, err := io.ReadAll(reader)
			if err != nil {
				t.Errorf("ReadAll() error = %v", err)
				return
			}
			if !bytes.Equal(got, data) {
				t.Errorf("ReadAll() got = %q, want %q", got, data)
			}
		}()
	}
	wg.Wait()
}

func TestMemoryBuffer_BinaryData(t *testing.T) {
	buf := NewMemoryBuffer()
	// Include null bytes, high bytes, and non-UTF8 sequences
	data := []byte{0x00, 0x01, 0xff, 0xfe, 0x80, 0x7f, 0x00}

	if _, err := buf.Write(data); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	buf.Close()

	reader := buf.NewReader(context.Background())
	defer reader.Close()

	got, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}
	if !bytes.Equal(got, data) {
		t.Fatalf("ReadAll() got = %v, want %v", got, data)
	}
}

func TestMemoryBuffer_BlocksUntilData(t *testing.T) {
	buf := NewMemoryBuffer()
	reader := buf.NewReader(context.Background())
	defer reader.Close()

	readDone := make(chan []byte, 1)
	go func() {
		got := make([]byte, 5)
		n, _ := reader.Read(got)
		readDone <- got[:n]
	}()

	// Should not receive anything yet
	select {
	case <-readDone:
		t.Fatal("Read should block until data is written")
	case <-time.After(50 * time.Millisecond):
		// expected
	}

	// Write data - should unblock reader
	if _, err := buf.Write([]byte("hello")); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	select {
	case got := <-readDone:
		if string(got) != "hello" {
			t.Fatalf("Read() got = %q, want %q", got, "hello")
		}
	case <-time.After(time.Second):
		t.Fatal("Read should have unblocked after Write")
	}

	buf.Close()
}

func TestMemoryBuffer_EOF_OnClose(t *testing.T) {
	buf := NewMemoryBuffer()
	data := []byte("test")
	if _, err := buf.Write(data); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	buf.Close()

	reader := buf.NewReader(context.Background())
	defer reader.Close()

	// First read should return data
	got := make([]byte, 10)
	n, err := reader.Read(got)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	if n != len(data) {
		t.Fatalf("Read() n = %d, want %d", n, len(data))
	}

	// Second read should return EOF
	n, err = reader.Read(got)
	if err != io.EOF {
		t.Fatalf("Read() error = %v, want io.EOF", err)
	}
	if n != 0 {
		t.Fatalf("Read() n = %d, want 0", n)
	}
}

func TestMemoryBuffer_CloseUnblocksReaders(t *testing.T) {
	buf := NewMemoryBuffer()
	reader := buf.NewReader(context.Background())
	defer reader.Close()

	readDone := make(chan error, 1)
	go func() {
		got := make([]byte, 1)
		_, err := reader.Read(got)
		readDone <- err
	}()

	// Should not receive anything yet
	select {
	case <-readDone:
		t.Fatal("Read should block on empty buffer")
	case <-time.After(50 * time.Millisecond):
		// expected
	}

	// Close should unblock reader with EOF
	buf.Close()

	select {
	case err := <-readDone:
		if err != io.EOF {
			t.Fatalf("Read() error = %v, want io.EOF", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Read should have unblocked after Close")
	}
}

func TestMemoryBuffer_MultipleWrites(t *testing.T) {
	buf := NewMemoryBuffer()
	reader := buf.NewReader(context.Background())
	defer reader.Close()

	// Write in chunks
	for _, chunk := range []string{"chunk1", "chunk2", "chunk3"} {
		if _, err := buf.Write([]byte(chunk)); err != nil {
			t.Fatalf("Write() error = %v", err)
		}
	}
	buf.Close()

	got, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}
	want := "chunk1chunk2chunk3"
	if string(got) != want {
		t.Fatalf("ReadAll() got = %q, want %q", got, want)
	}
}

func TestMemoryBuffer_ReadInSmallChunks(t *testing.T) {
	buf := NewMemoryBuffer()
	data := []byte("hello world this is a longer message")
	if _, err := buf.Write(data); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	buf.Close()

	reader := buf.NewReader(context.Background())
	defer reader.Close()

	var result []byte
	chunk := make([]byte, 5)

	for {
		n, err := reader.Read(chunk)
		result = append(result, chunk[:n]...)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Read() error = %v", err)
		}
	}

	if !bytes.Equal(result, data) {
		t.Fatalf("got = %q, want %q", result, data)
	}
}

func TestMemoryBuffer_Size(t *testing.T) {
	buf := NewMemoryBuffer()

	if buf.Size() != 0 {
		t.Fatalf("Size() = %d, want 0", buf.Size())
	}

	if _, err := buf.Write([]byte("hello")); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if buf.Size() != 5 {
		t.Fatalf("Size() = %d, want 5", buf.Size())
	}

	if _, err := buf.Write([]byte(" world")); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if buf.Size() != 11 {
		t.Fatalf("Size() = %d, want 11", buf.Size())
	}

	buf.Close()
}

func TestMemoryBuffer_DoubleClose(t *testing.T) {
	buf := NewMemoryBuffer()
	buf.Write([]byte("test"))

	// Should not panic
	buf.Close()
	buf.Close()
	buf.Close()

	if !buf.IsClosed() {
		t.Fatal("IsClosed() = false after Close()")
	}
}

func TestMemoryBuffer_IsClosed(t *testing.T) {
	buf := NewMemoryBuffer()

	if buf.IsClosed() {
		t.Fatal("IsClosed() = true for new buffer")
	}

	buf.Close()

	if !buf.IsClosed() {
		t.Fatal("IsClosed() = false after Close()")
	}
}

func TestMemoryBuffer_WriteToClosedBuffer(t *testing.T) {
	buf := NewMemoryBuffer()
	buf.Close()

	_, err := buf.Write([]byte("test"))
	if err != io.ErrClosedPipe {
		t.Fatalf("Write() error = %v, want io.ErrClosedPipe", err)
	}
}

func TestMemoryBufferReader_Offset(t *testing.T) {
	buf := NewMemoryBuffer()
	buf.Write([]byte("hello world"))
	buf.Close()

	reader := buf.NewReader(context.Background())
	defer reader.Close()

	if reader.Offset() != 0 {
		t.Fatalf("Offset() = %d, want 0", reader.Offset())
	}

	p := make([]byte, 5)
	reader.Read(p)

	if reader.Offset() != 5 {
		t.Fatalf("Offset() = %d, want 5", reader.Offset())
	}
}

func TestMemoryBufferReader_Close(t *testing.T) {
	buf := NewMemoryBuffer()
	reader := buf.NewReader(context.Background())

	// Close should return nil and stop the AfterFunc
	if err := reader.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	// Should be safe to close multiple times
	if err := reader.Close(); err != nil {
		t.Fatalf("Close() second call error = %v", err)
	}
}

func TestMemoryBufferReader_CloseUnblocksRead(t *testing.T) {
	// This tests the fix for the PR review comment:
	// Reader.Close() must unblock blocked Read() calls even without context cancellation
	buf := NewMemoryBuffer()
	reader := buf.NewReader(context.Background()) // No context cancellation

	done := make(chan error, 1)
	go func() {
		p := make([]byte, 10)
		_, err := reader.Read(p)
		done <- err
	}()

	// Should be blocked (no data, buffer not closed)
	select {
	case <-done:
		t.Fatal("Read should block on empty buffer")
	case <-time.After(50 * time.Millisecond):
	}

	// Close reader (not buffer!) - should unblock the Read
	reader.Close()

	select {
	case err := <-done:
		if err != ErrReaderClosed {
			t.Fatalf("Read() error = %v, want ErrReaderClosed", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Read should have unblocked after reader.Close()")
	}
}

func TestMemoryBufferReader_CloseIdempotent(t *testing.T) {
	buf := NewMemoryBuffer()
	reader := buf.NewReader(context.Background())

	// Multiple Close calls should be safe
	for i := 0; i < 5; i++ {
		if err := reader.Close(); err != nil {
			t.Fatalf("Close() call %d error = %v", i+1, err)
		}
	}
}

func TestMemoryBufferReader_ReadAfterClose(t *testing.T) {
	buf := NewMemoryBuffer()
	buf.Write([]byte("data"))

	reader := buf.NewReader(context.Background())
	reader.Close()

	// Read after Close should return ErrReaderClosed
	p := make([]byte, 10)
	_, err := reader.Read(p)
	if err != ErrReaderClosed {
		t.Fatalf("Read() error = %v, want ErrReaderClosed", err)
	}
}

func TestMemoryBufferReader_ContextCancellation(t *testing.T) {
	buf := NewMemoryBuffer()
	ctx, cancel := context.WithCancel(context.Background())

	reader := buf.NewReader(ctx)
	defer reader.Close()

	done := make(chan error, 1)
	go func() {
		p := make([]byte, 10)
		_, err := reader.Read(p)
		done <- err
	}()

	// Should be blocked
	select {
	case <-done:
		t.Fatal("Read should block on empty buffer")
	case <-time.After(50 * time.Millisecond):
	}

	// Cancel context - should unblock reader
	cancel()

	select {
	case err := <-done:
		if err != context.Canceled {
			t.Fatalf("Read() error = %v, want context.Canceled", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Read should have unblocked after context cancel")
	}
}

func TestMemoryBufferReader_ContextCancellation_WithData(t *testing.T) {
	buf := NewMemoryBuffer()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Write some data first
	buf.Write([]byte("hello"))

	reader := buf.NewReader(ctx)
	defer reader.Close()

	// Should be able to read the data even with cancellable context
	p := make([]byte, 5)
	n, err := reader.Read(p)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	if n != 5 {
		t.Fatalf("Read() n = %d, want 5", n)
	}
	if string(p) != "hello" {
		t.Fatalf("Read() got = %q, want %q", p, "hello")
	}
}

func TestMemoryBufferReader_ContextAlreadyCancelled(t *testing.T) {
	buf := NewMemoryBuffer()
	buf.Write([]byte("data"))

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel before creating reader

	reader := buf.NewReader(ctx)
	defer reader.Close()

	p := make([]byte, 10)
	_, err := reader.Read(p)
	if err != context.Canceled {
		t.Fatalf("Read() error = %v, want context.Canceled", err)
	}
}

func TestMemoryBufferReader_JobFinishesBeforeContextCancel(t *testing.T) {
	// This tests the normal case: job finishes, reader gets EOF, context never cancelled
	buf := NewMemoryBuffer()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cleanup, but shouldn't affect the test

	reader := buf.NewReader(ctx)
	defer reader.Close()

	// Write data and close (simulating job completion)
	buf.Write([]byte("output"))
	buf.Close()

	got, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}
	if string(got) != "output" {
		t.Fatalf("ReadAll() got = %q, want %q", got, "output")
	}
}

func TestMemoryBufferReader_ConcurrentReadersWithContexts(t *testing.T) {
	buf := NewMemoryBuffer()
	data := []byte("shared data for all readers")

	buf.Write(data)
	buf.Close()

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			reader := buf.NewReader(ctx)
			defer reader.Close()

			got, err := io.ReadAll(reader)
			if err != nil {
				t.Errorf("ReadAll() error = %v", err)
				return
			}
			if !bytes.Equal(got, data) {
				t.Errorf("ReadAll() got = %q, want %q", got, data)
			}
		}()
	}
	wg.Wait()
}

func TestMemoryBufferReader_StreamingWithContext(t *testing.T) {
	// Simulates gRPC streaming pattern
	buf := NewMemoryBuffer()
	ctx, cancel := context.WithCancel(context.Background())

	reader := buf.NewReader(ctx)
	defer reader.Close()

	received := make(chan []byte, 10)
	done := make(chan error, 1)

	go func() {
		p := make([]byte, 1024)
		for {
			n, err := reader.Read(p)
			if err != nil {
				done <- err
				return
			}
			chunk := make([]byte, n)
			copy(chunk, p[:n])
			received <- chunk
		}
	}()

	// Write chunks
	buf.Write([]byte("chunk1"))
	time.Sleep(10 * time.Millisecond)
	buf.Write([]byte("chunk2"))
	time.Sleep(10 * time.Millisecond)

	// Cancel context (simulating client disconnect)
	cancel()

	select {
	case err := <-done:
		if err != context.Canceled {
			t.Fatalf("Read() error = %v, want context.Canceled", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Reader should have exited after context cancel")
	}

	// Verify we received the chunks before cancel
	var got []byte
	for {
		select {
		case chunk := <-received:
			got = append(got, chunk...)
		default:
			goto done
		}
	}
done:
	if !bytes.HasPrefix(got, []byte("chunk1")) {
		t.Errorf("Expected to receive at least chunk1, got = %q", got)
	}
}

func TestMemoryBuffer_ReaderCreatedAfterClose(t *testing.T) {
	// Late joiner scenario: buffer is already closed with data
	buf := NewMemoryBuffer()
	data := []byte("historical output")

	buf.Write(data)
	buf.Close()

	// Reader created AFTER close should still get all data
	reader := buf.NewReader(context.Background())
	defer reader.Close()

	got, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}
	if !bytes.Equal(got, data) {
		t.Fatalf("ReadAll() got = %q, want %q", got, data)
	}
}

func TestMemoryBuffer_LargeData(t *testing.T) {
	// Test data larger than initial capacity (4096 bytes)
	buf := NewMemoryBuffer()

	// 10KB of data
	data := make([]byte, 10*1024)
	for i := range data {
		data[i] = byte(i % 256)
	}

	n, err := buf.Write(data)
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if n != len(data) {
		t.Fatalf("Write() n = %d, want %d", n, len(data))
	}

	if buf.Size() != int64(len(data)) {
		t.Fatalf("Size() = %d, want %d", buf.Size(), len(data))
	}

	buf.Close()

	reader := buf.NewReader(context.Background())
	defer reader.Close()

	got, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}
	if !bytes.Equal(got, data) {
		t.Fatal("Large data mismatch")
	}
}

func TestMemoryBuffer_EmptyWrite(t *testing.T) {
	buf := NewMemoryBuffer()

	// Empty write should succeed
	n, err := buf.Write([]byte{})
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if n != 0 {
		t.Fatalf("Write() n = %d, want 0", n)
	}

	if buf.Size() != 0 {
		t.Fatalf("Size() = %d, want 0", buf.Size())
	}

	// Write actual data after empty write
	buf.Write([]byte("data"))
	buf.Close()

	reader := buf.NewReader(context.Background())
	defer reader.Close()

	got, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}
	if string(got) != "data" {
		t.Fatalf("got = %q, want %q", got, "data")
	}
}

// === Tests for Bug Fixes ===

// TestMemoryBufferReader_CloseRaceStress tests the lost wakeup bug fix.
// It rapidly creates readers and closes them while reads are pending.
// Without the fix, some readers would block forever.
func TestMemoryBufferReader_CloseRaceStress(t *testing.T) {
	for i := 0; i < 100; i++ {
		buf := NewMemoryBuffer()
		reader := buf.NewReader(context.Background())

		done := make(chan error, 1)
		go func() {
			p := make([]byte, 10)
			_, err := reader.Read(p)
			done <- err
		}()

		// Give the goroutine time to enter the read loop
		// but not necessarily to reach Wait()
		time.Sleep(time.Microsecond * time.Duration(i%10))

		// Close reader - this must unblock the Read
		reader.Close()

		select {
		case err := <-done:
			if err != ErrReaderClosed && err != context.Canceled {
				t.Fatalf("iteration %d: Read() error = %v, want ErrReaderClosed", i, err)
			}
		case <-time.After(100 * time.Millisecond):
			t.Fatalf("iteration %d: Read blocked forever - lost wakeup bug!", i)
		}
	}
}

// TestMemoryBufferReader_OffsetThreadSafety tests that Offset() is thread-safe.
func TestMemoryBufferReader_OffsetThreadSafety(t *testing.T) {
	buf := NewMemoryBuffer()
	buf.Write(bytes.Repeat([]byte("x"), 10000))
	buf.Close()

	reader := buf.NewReader(context.Background())
	defer reader.Close()

	var wg sync.WaitGroup

	// Reader goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		p := make([]byte, 100)
		for {
			_, err := reader.Read(p)
			if err == io.EOF {
				return
			}
		}
	}()

	// Offset checker goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			offset := reader.Offset()
			if offset < 0 {
				t.Errorf("Offset() = %d, should never be negative", offset)
			}
			time.Sleep(time.Microsecond)
		}
	}()

	wg.Wait()
}

// TestMemoryBuffer_ConcurrentWriteAndRead tests concurrent write and read.
func TestMemoryBuffer_ConcurrentWriteAndRead(t *testing.T) {
	buf := NewMemoryBuffer()
	reader := buf.NewReader(context.Background())
	defer reader.Close()

	var wg sync.WaitGroup
	totalWrites := 1000
	writeSize := 100

	// Writer goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < totalWrites; i++ {
			buf.Write(bytes.Repeat([]byte{byte(i)}, writeSize))
		}
		buf.Close()
	}()

	// Reader goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		totalRead := 0
		p := make([]byte, 1024)
		for {
			n, err := reader.Read(p)
			totalRead += n
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Errorf("Read() error = %v", err)
				return
			}
		}
		expected := totalWrites * writeSize
		if totalRead != expected {
			t.Errorf("total read = %d, want %d", totalRead, expected)
		}
	}()

	wg.Wait()
}

// TestMemoryBuffer_PerReaderChannels tests that each reader has its own
// notification channel.
func TestMemoryBuffer_PerReaderChannels(t *testing.T) {
	buf := NewMemoryBuffer()

	// Create multiple readers
	numReaders := 10
	readers := make([]*MemoryBufferReader, numReaders)
	for i := 0; i < numReaders; i++ {
		readers[i] = buf.NewReader(context.Background())
	}

	// Write data - each reader should be notified independently
	buf.Write([]byte("hello"))

	// All readers should be able to read the data
	for i, reader := range readers {
		p := make([]byte, 10)
		n, err := reader.Read(p)
		if err != nil {
			t.Fatalf("reader %d: Read() error = %v", i, err)
		}
		if n != 5 || string(p[:n]) != "hello" {
			t.Fatalf("reader %d: Read() = %q, want %q", i, string(p[:n]), "hello")
		}
	}

	// Clean up
	for _, reader := range readers {
		reader.Close()
	}
	buf.Close()
}

// TestMemoryBuffer_ReaderUnregistration tests that readers are properly
// unregistered when closed.
func TestMemoryBuffer_ReaderUnregistration(t *testing.T) {
	buf := NewMemoryBuffer()

	// Create readers
	reader1 := buf.NewReader(context.Background())
	reader2 := buf.NewReader(context.Background())

	// Check both are registered
	buf.mu.Lock()
	if len(buf.readers) != 2 {
		t.Fatalf("expected 2 registered readers, got %d", len(buf.readers))
	}
	buf.mu.Unlock()

	// Close one reader
	reader1.Close()

	// Check only one remains
	buf.mu.Lock()
	if len(buf.readers) != 1 {
		t.Fatalf("expected 1 registered reader after close, got %d", len(buf.readers))
	}
	buf.mu.Unlock()

	// Close the other
	reader2.Close()

	// Check none remain
	buf.mu.Lock()
	if len(buf.readers) != 0 {
		t.Fatalf("expected 0 registered readers after close, got %d", len(buf.readers))
	}
	buf.mu.Unlock()

	buf.Close()
}

// TestMemoryBuffer_IndependentReaderNotification tests that readers wake
func TestMemoryBuffer_IndependentReaderNotification(t *testing.T) {
	buf := NewMemoryBuffer()

	// Create two readers
	reader1 := buf.NewReader(context.Background())
	reader2 := buf.NewReader(context.Background())
	defer reader1.Close()
	defer reader2.Close()

	// Write some data
	buf.Write([]byte("data1"))

	// Reader1 reads the data
	p := make([]byte, 10)
	n, _ := reader1.Read(p)
	if string(p[:n]) != "data1" {
		t.Fatalf("reader1: got %q, want %q", string(p[:n]), "data1")
	}

	// Reader1 is now at offset 5, reader2 at offset 0

	// Write more data
	buf.Write([]byte("data2"))

	// Reader2 should get all data from beginning
	all := make([]byte, 20)
	n1, _ := reader2.Read(all)
	if string(all[:n1]) != "data1data2" {
		t.Fatalf("reader2: got %q, want %q", string(all[:n1]), "data1data2")
	}

	// Reader1 should only get new data
	n2, _ := reader1.Read(p)
	if string(p[:n2]) != "data2" {
		t.Fatalf("reader1 second read: got %q, want %q", string(p[:n2]), "data2")
	}

	buf.Close()
}

// TestMemoryBuffer_CloseNotifiesAllReaders tests that Close() notifies
// all registered readers.
func TestMemoryBuffer_CloseNotifiesAllReaders(t *testing.T) {
	buf := NewMemoryBuffer()

	numReaders := 5
	var wg sync.WaitGroup

	for i := 0; i < numReaders; i++ {
		reader := buf.NewReader(context.Background())
		wg.Add(1)
		go func(r *MemoryBufferReader, id int) {
			defer wg.Done()
			defer r.Close()

			p := make([]byte, 10)
			_, err := r.Read(p)
			if err != io.EOF {
				t.Errorf("reader %d: expected EOF, got %v", id, err)
			}
		}(reader, i)
	}

	// Give readers time to block
	time.Sleep(50 * time.Millisecond)

	// Close should wake all readers
	buf.Close()

	// Wait for all readers with timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// All readers woke up
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for readers to wake up after Close()")
	}
}

// TestMemoryBuffer_ConcurrentWriters tests that multiple concurrent writers
// don't block each other (lock released before notifications).
func TestMemoryBuffer_ConcurrentWriters(t *testing.T) {
	buf := NewMemoryBuffer()
	defer buf.Close()

	// Create a reader to ensure notifications happen
	reader := buf.NewReader(context.Background())
	defer reader.Close()

	numWriters := 10
	writesPerWriter := 100
	var wg sync.WaitGroup

	start := time.Now()

	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			data := []byte(fmt.Sprintf("writer%d-", id))
			for j := 0; j < writesPerWriter; j++ {
				buf.Write(data)
			}
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)

	// Verify all data was written
	expectedSize := int64(numWriters * writesPerWriter * len("writer0-"))
	if buf.Size() != expectedSize {
		t.Errorf("Size() = %d, want %d", buf.Size(), expectedSize)
	}

	// With lock-before-send, this should complete quickly
	// With lock-during-send, writers would block each other
	if elapsed > 500*time.Millisecond {
		t.Logf("Warning: concurrent writes took %v (may indicate contention)", elapsed)
	}
}

// TestMemoryBuffer_ReaderCloseDuringWrite tests that a reader can close
// while Write() is sending notifications (safe after lock release).
func TestMemoryBuffer_ReaderCloseDuringWrite(t *testing.T) {
	buf := NewMemoryBuffer()

	// Create many readers
	numReaders := 100
	readers := make([]*MemoryBufferReader, numReaders)
	for i := 0; i < numReaders; i++ {
		readers[i] = buf.NewReader(context.Background())
	}

	var wg sync.WaitGroup

	// Writer goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			buf.Write([]byte("data"))
		}
	}()

	// Close readers concurrently while writes are happening
	for _, reader := range readers {
		wg.Add(1)
		go func(r *MemoryBufferReader) {
			defer wg.Done()
			time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
			r.Close()
		}(reader)
	}

	wg.Wait()
	buf.Close()

	// If we get here without panic/race, the test passes
}

// TestMemoryBuffer_WriteReadContention tests that readers can read
// while Write() is sending notifications (not blocked).
func TestMemoryBuffer_WriteReadContention(t *testing.T) {
	buf := NewMemoryBuffer()

	numReaders := 50
	var wg sync.WaitGroup

	// Start readers
	for i := 0; i < numReaders; i++ {
		reader := buf.NewReader(context.Background())
		wg.Add(1)
		go func(r *MemoryBufferReader) {
			defer wg.Done()
			defer r.Close()
			p := make([]byte, 1024)
			for {
				_, err := r.Read(p)
				if err == io.EOF {
					return
				}
				if err != nil {
					return
				}
			}
		}(reader)
	}

	// Writer - many small writes
	for i := 0; i < 1000; i++ {
		buf.Write([]byte("x"))
	}
	buf.Close()

	// Wait with timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Success
	case <-time.After(5 * time.Second):
		t.Fatal("timeout - possible deadlock from write/read contention")
	}
}
