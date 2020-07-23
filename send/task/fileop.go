package task

import (
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/Ghamster0/os-rq-fsender/pkg/hdfsfile"
	"github.com/Ghamster0/os-rq-fsender/pkg/sth"
	"github.com/Ghamster0/os-rq-fsender/send/entity"
	"github.com/colinmarc/hdfs"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type FileRef interface {
	io.Reader
	io.Seeker
	io.Closer
	Stat() (os.FileInfo, error)
}

type FileService struct {
	db *gorm.DB
}

func NewFileService(db *gorm.DB) *FileService {
	return &FileService{
		db: db,
	}
}

func (fs *FileService) SaveFileHDFS() string {
	return ""
}

func (fs *FileService) SaveFileLocal(path string, name string) (r sth.Result, err error) {
	meta := map[string]interface{}{
		"type":        "local",
		"path":        path,
		"origin_name": name,
	}
	var fr FileRef
	if fr, err = openFileLocal(path); err == nil {
		r, err = fs.saveFile(fr, meta)
	}
	return
}

func (fs *FileService) saveFile(fr FileRef, meta map[string]interface{}) (r sth.Result, err error) {
	var fileInfo os.FileInfo
	if fileInfo, err = fr.Stat(); err == nil {
		metaStr, _ := json.Marshal(meta)
		model := &entity.FileModel{
			Id:        uuid.NewV4().String(),
			Size:      fileInfo.Size(),
			FileMeta:  string(metaStr),
			Offset:    0,
			Success:   0,
			Fail:      0,
			Status:    entity.Idel,
			CreatedAt: time.Now(),
		}
		fs.db.Save(model)
		r = meta
		r["id"] = model.Id
		r["size"] = model.Size
	}
	return
}

func (fs *FileService) ListFiles(start int, limit int) ([]sth.Result, error) {
	var files []entity.FileModel
	if err := fs.db.Offset(start).Limit(limit).Find(&files).Error; err == nil {
		var res []sth.Result
		for i := range files {
			if r, err := files[i].Info(); err == nil {
				res = append(res, r)
			}
		}
		return res, nil
	} else {
		return nil, err
	}
}

func openFileHDFS(path string, host string) (fr FileRef, err error) {
	var client *hdfs.Client
	client, err = hdfs.New(host)
	if err == nil {
		var reader *hdfs.FileReader
		reader, err = client.Open(path)
		if err == nil {
			fr = &hdfsfile.WrappedFileReader{FileReader: reader}
		}
	}
	return
}

func openFileLocal(path string) (fr FileRef, err error) {
	fr, err = os.Open(path)
	return
}

func deleteFileLocal(path string) error {
	return os.Remove(path)
}
