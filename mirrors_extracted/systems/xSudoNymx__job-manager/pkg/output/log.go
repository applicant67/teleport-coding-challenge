package output

import (
	"errors"
	"fmt"
	"os"
	"sync"
)

// ErrLogClosed is returned when writing to a closed log.
var ErrLogClosed = errors.New("output log is closed")

type Log struct {
	// immutable
	file *os.File
	path string

	// mutable
	mu sync.Mutex

	// committedLen is total bytes written to disk.
	committedLen int64

	// Signal for subscribers to wake.
	notify chan struct{}

	// closed indicates the log is closed
	closed bool
}

type Config struct {
	// Path to file
	Path string
}

func New(cfg Config) (*Log, error) {
	file, err := os.Create(cfg.Path)
	if err != nil {
		return nil, fmt.Errorf("create output log: %w", err)
	}

	return &Log{
		file:   file,
		path:   cfg.Path,
		notify: make(chan struct{}),
	}, nil
}

// Append writes a framed record and notifies subscribers.
func (l *Log) Append(tag StreamTag, data []byte) (int, error) {
	frame := &Frame{Tag: tag, Data: data}
	buf, err := frame.MarshalBinary()
	if err != nil {
		return 0, err
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.closed {
		return 0, ErrLogClosed
	}

	// TODO: Handle partial writes. If n > 0 && err != nil.
	//  A partial frame may be written to disk, causing truncated replay.
	n, err := l.file.Write(buf)
	if err != nil {
		return n, fmt.Errorf("write frame: %w", err)
	}

	l.committedLen += int64(n)

	// Notify all waiting subscribers by closing the channel
	close(l.notify)
	// Create new channel for next notification
	l.notify = make(chan struct{})

	return n, nil
}

// Close closes the log and signals subscribers to finish.
func (l *Log) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.closed {
		return nil
	}

	l.closed = true
	// Wake any waiting subscribers
	close(l.notify)

	return l.file.Close()
}

// getState returns a snapshot of the current output log state.
func (l *Log) getState() (notify <-chan struct{}, committedLen int64, closed bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.notify, l.committedLen, l.closed
}

// Subscribe creates a new subscription and returns a subscriber.
func (l *Log) Subscribe() (*Subscriber, error) {
	// Open file handle for this subscriber
	file, err := os.Open(l.path)
	if err != nil {
		return nil, fmt.Errorf("open log for subscription: %w", err)
	}

	return &Subscriber{
		log:    l,
		file:   file,
		offset: 0,
	}, nil
}
