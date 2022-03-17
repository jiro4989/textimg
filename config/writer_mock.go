package config

import (
	"errors"
	"io"
)

type MockWriter struct {
	writeErr bool
	closeErr bool
}

func NewMockWriter(w, c bool) io.WriteCloser {
	return &MockWriter{
		writeErr: w,
		closeErr: c,
	}
}

func (m *MockWriter) Write(p []byte) (n int, err error) {
	if m.writeErr {
		return -1, errors.New("write error")
	}
	return 0, nil
}

func (m *MockWriter) Close() error {
	if m.closeErr {
		return errors.New("close error")
	}
	return nil
}
