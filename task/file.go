package task

import (
	"io"
	"os"

	"github.com/Ghamster0/os-rq-fsender/conf"
)

// PlainFile TODO
type PlainFile interface {
	Open() (VisableReader, error)
	FileMeta() Result
}

// VisableReader TODO
type VisableReader interface {
	io.Reader
	io.Seeker
	Size() (int64, error)
}

// OSFile VisableReader implemention
type OSFile struct {
	*os.File
}

// Size TODO
func (f *OSFile) Size() (size int64, err error) {
	var fileInfo os.FileInfo
	fileInfo, err = f.Stat()
	if err != nil {
		return
	}
	size = fileInfo.Size()
	return
}

// LocalFile TODO
type LocalFile struct {
	Fid string
}

// Open return io.Reader
func (lf *LocalFile) Open() (osf VisableReader, err error) {
	var f *os.File
	f, err = os.Open(conf.FileStore + lf.Fid)
	if err == nil {
		osf = &OSFile{f}
	}
	return
}

// FileMeta TODO
func (lf *LocalFile) FileMeta() Result {
	res := Result{
		"fid": lf.Fid,
	}
	return res
}

// LocalFileLoader TODO
func LocalFileLoader(fnames *[]string) []PlainFile {
	l := len(*fnames)
	var files = make([]PlainFile, l)
	for i := 0; i < l; i++ {
		files[i] = &LocalFile{Fid: (*fnames)[i]}
	}
	return files
}
