package task

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/Ghamster0/os-rq-fsender/send/entity"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type CleanService struct {
	db            *gorm.DB
	uploadPath    string
	uploadLogPath string
}

func NewCleanService(db *gorm.DB, conf *viper.Viper) *CleanService {
	return &CleanService{
		db:            db,
		uploadPath:    conf.GetString("upload.file"),
		uploadLogPath: conf.GetString("upload.log"),
	}
}

func CleanServiceServ(lifecycle fx.Lifecycle, cs *CleanService) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go cs.StartCleanService()
				return nil
			},
		},
	)
}

func (cs *CleanService) StartCleanService() {
	for {
		select {
		case <-time.After(time.Minute * 10):
			cs.cleanUpload()
			cs.cleanLog()
		}
	}
}

func (cs *CleanService) cleanUpload() {
	var fms []entity.FileModel
	if err := cs.db.Where("status = ? AND created_at < ?", entity.Idel, time.Now().Add(-time.Hour*12)).Find(&fms).Error; err != nil {
		return
	}
	for i := range fms {
		fm := fms[i]
		var fileMeta entity.FileMeta
		if err := json.Unmarshal([]byte(fm.FileMeta), &fileMeta); err == nil {
			if t := fileMeta["type"].(string); t == "local" {
				// remove local file
				os.Remove(fileMeta["path"].(string))
			}
		}
		cs.db.Delete(&fm)
	}
}

func (cs *CleanService) cleanLog() {
	fileInfos, err := ioutil.ReadDir(cs.uploadLogPath)
	if err != nil {
		return
	}
	now := time.Now()
	for _, info := range fileInfos {
		if diff := now.Sub(info.ModTime()); diff > time.Hour*24*7 {
			os.Remove(cs.uploadLogPath + info.Name())
		}
	}
}
