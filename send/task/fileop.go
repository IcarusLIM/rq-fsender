package task

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
	"time"

	"github.com/Ghamster0/os-rq-fsender/pkg/dto"
	"github.com/Ghamster0/os-rq-fsender/pkg/hdfsfile"
	"github.com/Ghamster0/os-rq-fsender/pkg/sth"
	"github.com/Ghamster0/os-rq-fsender/send/entity"
	"github.com/colinmarc/hdfs/v2"
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

func (fs *FileService) SaveFileHDFS(path string, hdfsConf *dto.HDFSConfig) (r sth.Result, err error) {
	pathArr := strings.Split(path, "/")
	confStr, _ := json.Marshal(hdfsConf)
	meta := map[string]interface{}{
		"type":        "hdfs",
		"path":        path,
		"origin_name": pathArr[len(pathArr)-1],
		"hdfs":        string(confStr),
	}
	var fr FileRef
	if fr, err = openFileHDFS(path, hdfsConf); err == nil {
		r, err = fs.saveFile(fr, meta)
	}
	return
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
			UpdateAt:  time.Now(),
		}
		fs.db.Save(model)
		r = meta
		r["id"] = model.Id
		r["size"] = model.Size
	}
	fr.Close()
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

func openFileHDFS(path string, hdfsConf *dto.HDFSConfig) (fr FileRef, err error) {
	err = errors.New("Empty NameNodes")
	for _, nameNode := range hdfsConf.NameNodes {
		clientOptions := hdfs.ClientOptions{
			Addresses: []string{nameNode},
			User:      hdfsConf.User,
		}
		var client *hdfs.Client
		client, err = hdfs.NewClient(clientOptions)
		if err != nil {
			continue
		}
		var f *hdfs.FileReader
		f, err = client.Open(path)
		if err != nil {
			client.Close()
			continue
		}
		fr = &hdfsfile.WrappedFileReader{FileReader: f, Client: client}
		return
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
