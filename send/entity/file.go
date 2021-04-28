package entity

import (
	"encoding/json"
	"time"

	"github.com/Ghamster0/os-rq-fsender/pkg/sth"
)

type FileStatus string

const (
	Fail       FileStatus = "fail"
	Cancel     FileStatus = "cancel"
	Idel       FileStatus = "idel"
	Waitting   FileStatus = "waitting"
	Processing FileStatus = "processing"
	Paused     FileStatus = "paused"
	Finish     FileStatus = "finish"
)

type FileModel struct {
	Id        string `gorm:"type:varchar(36);primary_key"`
	BatchId   string
	FileMeta  string
	Size      int64
	Offset    int64
	Success   int32
	Fail      int32
	Status    FileStatus
	Err       string
	CreatedAt time.Time
	UpdateAt  time.Time
}

type FileMeta map[string]interface{}

func (meta FileMeta) GetType() string {
	return meta["type"].(string)
}

func (meta FileMeta) GetPath() string {
	return meta["path"].(string)
}

func (fm *FileModel) Info() (r sth.Result, err error) {
	err = json.Unmarshal([]byte(fm.FileMeta), &r)
	if err == nil {
		r["id"] = fm.Id
		r["success"] = fm.Success
		r["fail"] = fm.Fail
		r["offset"] = fm.Offset
		r["size"] = fm.Size
		r["status"] = fm.Status
	}
	return
}
