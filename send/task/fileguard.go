package task

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/Ghamster0/os-rq-fsender/pkg/sth"
	"github.com/Ghamster0/os-rq-fsender/send/entity"
	"github.com/jinzhu/gorm"
)

type Guard struct {
	id     string
	db     *gorm.DB
	model  *entity.FileModel
	meta   entity.FileMeta
	fr     FileRef
	reader *bufio.Reader
	offset int64
}

// not thread safe
func GuardOpen(fid string, db *gorm.DB) (g *Guard, success int, fail int, err error) {
	var fileModel entity.FileModel
	var fileMeta entity.FileMeta

	db.Where("id = ?", fid).First(&fileModel)
	err = json.Unmarshal([]byte(fileModel.FileMeta), &fileMeta)
	if err == nil {
		var fileRef FileRef
		if t := fileMeta["type"].(string); t == "local" {
			fileRef, err = openFileLocal(fileMeta["path"].(string))
		} else {
			fileRef, err = openFileHDFS(fileMeta["path"].(string), fileMeta["host"].(string))
		}
		if err == nil {
			success = fileModel.Success
			fail = fileModel.Fail
			offset := fileModel.Offset
			fileRef.Seek(offset, os.SEEK_SET)
			reader := bufio.NewReader(fileRef)
			g = &Guard{
				id:     fid,
				db:     db,
				model:  &fileModel,
				meta:   fileMeta,
				fr:     fileRef,
				reader: reader,
				offset: offset,
			}
		}
	}
	return
}

func (guard *Guard) GuardClose(closeAt entity.FileStatus) {
	guard.fr.Close()
	meta, _ := guard.Info()
	logger.Info("File", meta["path"], "close - at", closeAt)
	guard.db.Model(guard.model).Updates(entity.FileModel{Status: closeAt})

	// On finish, remove file
	if closeAt == entity.Finish {
		if guard.meta.GetType() == "local" {
			if err := deleteFileLocal(guard.meta.GetPath()); err != nil {
				logger.Warning("Can't remove file: ", guard.meta.GetPath())
			}
		}
	}
}

func (guard *Guard) ReadLine() (string, error) {
	line, err := guard.reader.ReadBytes('\n')
	guard.offset += int64(len(line))
	return string(line), err
}

func (guard *Guard) Update(success int, fail int) {
	guard.db.Model(guard.model).Updates(entity.FileModel{Success: success, Fail: fail, Offset: guard.offset})
}

func (guard *Guard) Cancel() {
	guard.GuardClose(entity.Cancel)
}

func (guard *Guard) Info() (meta sth.Result, err error) {
	if meta, err = guard.model.Info(); err == nil {
		meta["offset"] = guard.offset
		meta["size"] = guard.model.Size
	}
	return
}
