package router

import (
	"github.com/Ghamster0/os-rq-fsender/app/controller"
	"github.com/gin-gonic/gin"
)

func BatchRouter(root *gin.RouterGroup, ctl *controller.BatchController) {
	g := root.Group("/batch/")
	g.POST("", ctl.CreateBatch)
	g.GET("", ctl.ListBatch)
	g.GET("/:id", ctl.GetBatch)
}
