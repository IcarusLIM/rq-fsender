package controller

import (
	"net/http"

	"github.com/Ghamster0/os-rq-fsender/pkg/dto"
	"github.com/Ghamster0/os-rq-fsender/pkg/sth"
	"github.com/Ghamster0/os-rq-fsender/send/task"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
)

type BatchController struct {
	db   *gorm.DB
	conf *viper.Viper
	bs   *task.BatchService
}

func NewBatchController(db *gorm.DB, conf *viper.Viper, bs *task.BatchService) *BatchController {
	return &BatchController{
		db:   db,
		conf: conf,
		bs:   bs,
	}
}

func (ctl *BatchController) CreateBatch(c *gin.Context) {
	var batchReq *dto.BatchReq = &dto.BatchReq{}
	var err error
	var res sth.Result
	if err = c.ShouldBindJSON(batchReq); err == nil {
		if batchReq.Id == "" {
			batchReq.Id = uuid.NewV1().String()
		}
		res, err = ctl.bs.CreateBatch(batchReq)
	}
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"res": true, "data": res})
	} else {
		c.JSON(http.StatusOK, gin.H{"res": false, "err": err.Error()})
	}
}

func (ctl *BatchController) ListBatch(c *gin.Context) {
	if res, err := ctl.bs.ListBatch(0, 10); err == nil {
		c.JSON(http.StatusOK, gin.H{"res": true, "data": res})
	} else {
		c.JSON(http.StatusOK, gin.H{"res": false, "err": err.Error()})
	}
}

func (ctl *BatchController) GetBatch(c *gin.Context) {
	if res, err := ctl.bs.GetBatch(c.Param("id")); err == nil {
		c.JSON(http.StatusOK, gin.H{"res": true, "data": res})
	} else {
		c.JSON(http.StatusOK, gin.H{"res": false, "err": err.Error()})
	}
}
