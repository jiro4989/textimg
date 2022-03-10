package config

import (
	"io"
	"os"
)

type MockWriter struct{}

func (f *MockWriter) NewWriter(path string) (io.WriteCloser, error) {
	return os.Create(path)
}
