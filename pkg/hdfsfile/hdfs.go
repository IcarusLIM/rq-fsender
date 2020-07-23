package hdfsfile

import (
	"os"

	"github.com/colinmarc/hdfs"
)

type WrappedFileReader struct {
	*hdfs.FileReader
}

func (w *WrappedFileReader) Stat() (os.FileInfo, error) {
	return w.FileReader.Stat(), nil
}
