package router

import (
	"github.com/Ghamster0/os-rq-fsender/app/controller"
	"github.com/Ghamster0/os-rq-fsender/app/ctx"
	"github.com/gin-gonic/gin"
)

// IRoutersHTTPFunc TODO
type IRoutersHTTPFunc func(string, ...gin.HandlerFunc) gin.IRoutes

// InitRouter TODO
func InitRouter(g gin.IRouter, ctx *ctx.ApplicationContext) {

	routers := []struct {
		HTTPFunc IRoutersHTTPFunc
		Path     string
		F        controller.CtrlFunc
	}{
		{g.POST, "/upload", controller.Upload},
		{g.GET, "/tasks", controller.ListTask},
		{g.POST, "/task", controller.SendFromLocal},
		{g.GET, "/task/:taskID", controller.TaskStatus},
		{g.DELETE, "/task/:taskID", controller.CancelTask},
	}

	wp := &controller.HandlerWrapper{App: ctx}
	for _, r := range routers {
		r.HTTPFunc(r.Path, wp.Wrap(r.F))
	}

}
