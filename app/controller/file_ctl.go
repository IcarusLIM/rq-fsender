package controller

import (
	"io"
	"net/http"
	"os"

	"github.com/Ghamster0/os-rq-fsender/pkg/dto"
	"github.com/Ghamster0/os-rq-fsender/pkg/sth"
	"github.com/Ghamster0/os-rq-fsender/send/entity"
	"github.com/Ghamster0/os-rq-fsender/send/task"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
)

type FileController struct {
	db          *gorm.DB
	conf        *viper.Viper
	fileService *task.FileService
	box         *task.TaskBox
}

func NewSendController(db *gorm.DB, conf *viper.Viper, fs *task.FileService, box *task.TaskBox) *FileController {
	return &FileController{
		db:          db,
		conf:        conf,
		fileService: fs,
		box:         box,
	}
}

func (ctl *FileController) Upload(c *gin.Context) {
	file, _ := c.FormFile("file")
	originName := file.Filename
	name := "upload_" + uuid.NewV4().String()
	filePath := ctl.conf.GetString("upload.file") + name
	c.SaveUploadedFile(file, filePath)
	res, err := ctl.fileService.SaveFileLocal(filePath, originName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"res": false, "err": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"res": true, "data": res})
	}
}

func (ctl *FileController) UploadHDFS(c *gin.Context) {
	var up *dto.UploadHDFS = &dto.UploadHDFS{}
	var err error
	var res sth.Result
	if err = c.ShouldBindJSON(up); err == nil {
		res, err = ctl.fileService.SaveFileHDFS(up.Path, &up.HDFSConf)
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"res": false, "err": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"res": true, "data": res})
	}
}

func (ctl *FileController) List(c *gin.Context) {
	if res, err := ctl.fileService.ListFiles(0, 10); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"res": false, "err": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"res": true, "data": res})
	}
}

func (ctl *FileController) Pause(c *gin.Context) {
	if err := ctl.box.PauseTask(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"res": false, "err": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"res": true})
	}
}

func (ctl *FileController) Resume(c *gin.Context) {
	if err := ctl.db.Where("id = ? AND status IN ( ? )", c.Param("id"), []entity.FileStatus{entity.Paused, entity.Fail}).Model(&entity.FileModel{}).Update("status", entity.Waitting).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"res": false, "err": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"res": true})
	}
}

func (ctl *FileController) GetLog(c *gin.Context) {
	filePath := ctl.conf.GetString("upload.log") + c.Param("id") + ".log"
	if file, err := os.Open(filePath); err == nil {
		defer file.Close()

		c.Writer.Header().Add("Content-type", "application/octet-stream")
		c.Writer.Header().Add("Content-Disposition", "attachment; filename=\"upload.log\"")
		if _, err = io.Copy(c.Writer, file); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"err": err.Error()})
		}
	} else {
		c.JSON(http.StatusNotFound, gin.H{"err": err.Error()})
	}
}
