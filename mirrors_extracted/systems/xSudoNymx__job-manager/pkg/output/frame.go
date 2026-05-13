package output

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

var (
	ErrPayloadTooLarge = errors.New("payload exceeds maximum size")
	ErrInvalidTag      = errors.New("invalid stream tag")
)

const (
	// MaxPayloadSize is the maximum payload per frame (32 KiB).
	MaxPayloadSize = 32 * 1024

	// FrameHeaderSize is tag (1 byte) + length (4 bytes).
	FrameHeaderSize = 5
)

type StreamTag byte

const (
	TagUnspecified StreamTag = 0x00
	TagStdout      StreamTag = 0x01
	TagStderr      StreamTag = 0x02
)

func (t StreamTag) String() string {
	switch t {
	case TagStdout:
		return "stdout"
	case TagStderr:
		return "stderr"
	default:
		return "unknown"
	}
}

// Frame represents a single output record.
//
// Wire format:
//
//	┌─────┬────────┬─────────┐
//	│ Tag │ Length │ Payload │
//	│ 1B  │ 4B LE  │ N bytes │
//	└─────┴────────┴─────────┘
type Frame struct {
	Tag  StreamTag
	Data []byte
}

// Size returns the total size of the frame.
func (f *Frame) Size() int {
	return FrameHeaderSize + len(f.Data)
}

// Validate checks if the frame is valid for serialization.
func (f *Frame) Validate() error {
	if f.Tag != TagStdout && f.Tag != TagStderr {
		return ErrInvalidTag
	}
	if len(f.Data) > MaxPayloadSize {
		return ErrPayloadTooLarge
	}
	return nil
}

// MarshalBinary encodes the bytes to frame format.
func (f *Frame) MarshalBinary() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}

	buf := make([]byte, FrameHeaderSize+len(f.Data))
	buf[0] = byte(f.Tag)
	binary.LittleEndian.PutUint32(buf[1:5], uint32(len(f.Data)))
	copy(buf[FrameHeaderSize:], f.Data)

	return buf, nil
}

// ReadFrame reads a single frame from the reader.
// Returns io.EOF if no more data is available.
func ReadFrame(r io.Reader) (*Frame, error) {
	// Read header
	var header [FrameHeaderSize]byte
	if _, err := io.ReadFull(r, header[:]); err != nil {
		return nil, err // Includes io.EOF
	}

	tag := StreamTag(header[0])
	if tag != TagStdout && tag != TagStderr {
		return nil, ErrInvalidTag
	}

	length := binary.LittleEndian.Uint32(header[1:5])
	if length > MaxPayloadSize {
		return nil, ErrPayloadTooLarge
	}

	// Read payload
	payload := make([]byte, length)
	if length > 0 {
		if _, err := io.ReadFull(r, payload); err != nil {
			return nil, fmt.Errorf("failed to read payload: %w", err)
		}
	}

	return &Frame{Tag: tag, Data: payload}, nil
}
