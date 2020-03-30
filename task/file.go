package task

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/Ghamster0/os-rq-fsender/conf"
	"github.com/go-redis/redis/v7"
)

// PlainFile TODO
type PlainFile interface {
	Open() (VisableReader, error)
	GetInfo() Result
	SetFlag(*Result)
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
	Fid    string
	Meta   *LocalFileMete
	status *Result
}

// LocalFileMete TODO
type LocalFileMete struct {
	Fid   string
	Name  string
	Ftype string
	Size  int64
	Time  string
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

// GetInfo TODO
func (lf *LocalFile) GetInfo() Result {
	res := Result{
		"fid":    lf.Fid,
		"status": lf.status,
	}
	return res
}

// SetFlag TODO
func (lf *LocalFile) SetFlag(r *Result) {
	lf.status = r
}

// LocalFileLoader TODO
func LocalFileLoader(fnames *[]string, rclient *redis.Client) []PlainFile {
	l := len(*fnames)
	var files = make([]PlainFile, l)
	for i := 0; i < l; i++ {
		var meta LocalFileMete
		t := rclient.HGet(conf.FilesKey, (*fnames)[i]).Val()
		json.Unmarshal([]byte(t), &meta)
		files[i] = &LocalFile{Fid: (*fnames)[i], Meta: &meta}
	}
	return files
}
