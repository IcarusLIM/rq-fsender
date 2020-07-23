package controller

import (
	"net/http"

	"github.com/Ghamster0/os-rq-fsender/send/task"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
)

type SendController struct {
	db          *gorm.DB
	conf        *viper.Viper
	fileService *task.FileService
}

func NewSendController(db *gorm.DB, conf *viper.Viper, fs *task.FileService) *SendController {
	return &SendController{
		db:          db,
		conf:        conf,
		fileService: fs,
	}
}

func (ctl *SendController) Upload(c *gin.Context) {
	file, _ := c.FormFile("file")
	originName := file.Filename
	name := "upload_" + uuid.NewV4().String()
	filePath := ctl.conf.GetString("upload.path") + name
	c.SaveUploadedFile(file, filePath)
	res, err := ctl.fileService.SaveFileLocal(filePath, originName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"res": true, "err": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"res": true, "data": res})
	}
}

func (ctl *SendController) List(c *gin.Context) {
	if res, err := ctl.fileService.ListFiles(0, 10); err != nil {
		c.JSON(http.StatusOK, gin.H{"res": true, "err": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"res": true, "data": res})
	}
}
