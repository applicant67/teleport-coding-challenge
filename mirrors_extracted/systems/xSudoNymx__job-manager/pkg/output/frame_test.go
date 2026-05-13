package output

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

func TestStreamTag_String(t *testing.T) {
	tests := []struct {
		tag  StreamTag
		want string
	}{
		{TagStdout, "stdout"},
		{TagStderr, "stderr"},
		{TagUnspecified, "unknown"},
	}

	for _, tt := range tests {
		if got := tt.tag.String(); got != tt.want {
			t.Errorf("StreamTag(%d).String() = %q, want %q", tt.tag, got, tt.want)
		}
	}
}

func TestFrame_MarshalBinary(t *testing.T) {
	tests := []struct {
		name string
		tag  StreamTag
		data []byte
	}{
		{"stdout empty", TagStdout, []byte{}},
		{"stderr empty", TagStderr, []byte{}},
		{"stdout text", TagStdout, []byte("hello world")},
		{"stderr text", TagStderr, []byte("error!")},
		{"binary data", TagStdout, []byte{0x00, 0xFF, 0x01, 0xFE}},
		{"max size", TagStdout, make([]byte, MaxPayloadSize)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frame := &Frame{Tag: tt.tag, Data: tt.data}

			encoded, err := frame.MarshalBinary()
			if err != nil {
				t.Fatalf("MarshalBinary() error = %v", err)
			}

			if len(encoded) != frame.Size() {
				t.Errorf("encoded size = %d, want %d", len(encoded), frame.Size())
			}
		})
	}
}

func TestFrame_Validate(t *testing.T) {
	tests := []struct {
		name    string
		frame   Frame
		wantErr error
	}{
		{"valid stdout", Frame{TagStdout, []byte("test")}, nil},
		{"valid stderr", Frame{TagStderr, []byte("test")}, nil},
		{"invalid tag", Frame{TagUnspecified, []byte("test")}, ErrInvalidTag},
		{"too large", Frame{TagStdout, make([]byte, MaxPayloadSize+1)}, ErrPayloadTooLarge},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.frame.Validate()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Validate() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestReadFrame(t *testing.T) {
	// Build frames manually using MarshalBinary
	frames := []struct {
		tag  StreamTag
		data []byte
	}{
		{TagStdout, []byte("first")},
		{TagStderr, []byte("second")},
		{TagStdout, []byte("third")},
	}

	var buf bytes.Buffer
	for _, f := range frames {
		frame := &Frame{Tag: f.tag, Data: f.data}
		encoded, err := frame.MarshalBinary()
		if err != nil {
			t.Fatalf("MarshalBinary() error = %v", err)
		}
		buf.Write(encoded)
	}

	// Read them back
	reader := bytes.NewReader(buf.Bytes())
	for i, want := range frames {
		got, err := ReadFrame(reader)
		if err != nil {
			t.Fatalf("ReadFrame(%d) error = %v", i, err)
		}
		if got.Tag != want.tag {
			t.Errorf("frame[%d].Tag = %v, want %v", i, got.Tag, want.tag)
		}
		if !bytes.Equal(got.Data, want.data) {
			t.Errorf("frame[%d].Data = %q, want %q", i, got.Data, want.data)
		}
	}

	// Next read should return EOF
	_, err := ReadFrame(reader)
	if err != io.EOF {
		t.Errorf("ReadFrame() at end = %v, want io.EOF", err)
	}
}

func TestReadFrame_Errors(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"too short", []byte{0x01, 0x02}},
		{"invalid tag", []byte{0x00, 0x00, 0x00, 0x00, 0x00}},
		{"truncated payload", []byte{0x01, 0x05, 0x00, 0x00, 0x00, 0xAA}},
		{"payload too large", []byte{0x01, 0x01, 0x00, 0x80, 0x00}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bytes.NewReader(tt.data)
			if _, err := ReadFrame(reader); err == nil {
				t.Error("ReadFrame() should return error")
			}
		})
	}
}

func TestFrame_Size(t *testing.T) {
	f := &Frame{Tag: TagStdout, Data: []byte("hello")}
	want := FrameHeaderSize + 5
	if got := f.Size(); got != want {
		t.Errorf("Size() = %d, want %d", got, want)
	}
}
