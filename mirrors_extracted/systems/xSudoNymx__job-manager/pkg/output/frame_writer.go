package output

// FrameWriter wraps a Log and tags all writes with a specific tag.
type FrameWriter struct {
	log *Log
	tag StreamTag
}

// NewFrameWriter creates a writer that tags output with the given tag.
func NewFrameWriter(log *Log, tag StreamTag) *FrameWriter {
	return &FrameWriter{log: log, tag: tag}
}

func (w *FrameWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	written := 0
	remaining := p

	for len(remaining) > 0 {
		chunk := remaining
		if len(chunk) > MaxPayloadSize {
			chunk = remaining[:MaxPayloadSize]
		}

		if _, err := w.log.Append(w.tag, chunk); err != nil {
			return written, err
		}

		written += len(chunk)
		remaining = remaining[len(chunk):]
	}

	return written, nil
}
