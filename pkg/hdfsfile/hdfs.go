package hdfsfile

import (
	"os"

	"github.com/colinmarc/hdfs/v2"
)

type WrappedFileReader struct {
	*hdfs.FileReader
	*hdfs.Client
}

func (w *WrappedFileReader) Stat() (os.FileInfo, error) {
	return w.FileReader.Stat(), nil
}

func (w *WrappedFileReader) Close() error {
	if err := w.FileReader.Close(); err != nil {
		return err
	}
	if err := w.Client.Close(); err != nil {
		return err
	}
	return nil
}
