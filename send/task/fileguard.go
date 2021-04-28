package task

import (
	"bufio"
	"encoding/json"
	"os"
	"sync"

	"github.com/Ghamster0/os-rq-fsender/pkg/dto"
	"github.com/Ghamster0/os-rq-fsender/send/entity"
	"github.com/jinzhu/gorm"
)

type Guard struct {
	db     *gorm.DB
	meta   entity.FileMeta
	fr     FileRef
	reader *bufio.Reader
	offset int64
	lock   *sync.RWMutex
}

// not thread safe
func GuardOpen(fileMeta entity.FileMeta, offset int64, db *gorm.DB) (g *Guard, err error) {
	var fileRef FileRef
	if t := fileMeta["type"].(string); t == "local" {
		fileRef, err = openFileLocal(fileMeta["path"].(string))
	} else {
		var hdfsConf *dto.HDFSConfig = &dto.HDFSConfig{}
		if err = json.Unmarshal([]byte(fileMeta["hdfs"].(string)), hdfsConf); err != nil {
			logger.Error(err)
		}
		fileRef, err = openFileHDFS(fileMeta["path"].(string), hdfsConf)
	}
	if err == nil {
		fileRef.Seek(offset, os.SEEK_SET)
		reader := bufio.NewReader(fileRef)
		g = &Guard{
			db:     db,
			meta:   fileMeta,
			fr:     fileRef,
			reader: reader,
			offset: offset,
			lock:   &sync.RWMutex{},
		}
	}
	return
}

func (guard *Guard) GuardClose(closeAt entity.FileStatus) {
	guard.fr.Close()
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
	guard.lock.Lock()
	defer guard.lock.Unlock()
	line, err := guard.reader.ReadBytes('\n')
	guard.offset += int64(len(line))
	return string(line), err
}

func (guard *Guard) Cancel() {
	guard.GuardClose(entity.Cancel)
}

func (guard *Guard) GuardOffset() int64 {
	return guard.offset
}
