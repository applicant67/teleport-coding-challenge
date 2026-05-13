package lib

import (
	"bytes"
	"sync"
)

// Buffer is a goroutine safe bytes.Buffer
type Buffer struct {
	sync.RWMutex
	buffer bytes.Buffer
}

// Write appends the contents of p to the buffer, growing the buffer as needed.
// It returns the number of bytes written.
func (b *Buffer) Write(p []byte) (n int, err error) {
	b.Lock()
	defer b.Unlock()
	return b.buffer.Write(p)
}

// String returns the contents of the unread portion of the buffer as a string.
// If the Buffer is a nil pointer, it returns "<nil>".
func (b *Buffer) String() string {
	b.RLock()
	defer b.RUnlock()
	return b.buffer.String()
}
