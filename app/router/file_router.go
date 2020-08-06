package router

import (
	"github.com/Ghamster0/os-rq-fsender/app/controller"
	"github.com/gin-gonic/gin"
)

func FileRouter(root *gin.RouterGroup, ctl *controller.FileController) {
	g := root.Group("/file/")
	g.POST("/upload", ctl.Upload)
	g.POST("/upload-hdfs", ctl.UploadHDFS)
	g.GET("", ctl.List)
	g.PUT("/:id/pause", ctl.Pause)
	g.PUT("/:id/resume", ctl.Resume)
}
