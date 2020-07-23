package router

import (
	"github.com/Ghamster0/os-rq-fsender/app/controller"
	"github.com/gin-gonic/gin"
)

func SendRouter(root *gin.RouterGroup, ctl *controller.SendController) {
	g := root.Group("/upload/")
	g.POST("", ctl.Upload)
	g.GET("", ctl.List)
}
