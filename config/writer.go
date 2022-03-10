package config

import (
	"io"
	"os"
)

type WriterInterface interface {
	NewWriter(path string) (io.WriteCloser, error)
}

type FileWriter struct{}

func (f *FileWriter) NewWriter(path string) (io.WriteCloser, error) {
	return os.Create(path)
}
