package task

import (
	"io"
	"os"
)

// ReadOnlyFile TODO
type ReadOnlyFile interface {
	Open() (*io.Reader, error)
}

// LocalFile TODO
type LocalFile struct {
	Path string
}

// Open TODO
func (f *LocalFile) Open() (file *os.File, err error) {
	file, err = os.Open(f.Path)
	return
}
