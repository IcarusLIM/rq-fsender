package controller

import (
	"github.com/Ghamster0/os-rq-fsender/app/ctx"
	"github.com/gin-gonic/gin"
)

// CtrlFunc TODO
type CtrlFunc func(*gin.Context, *ctx.ApplicationContext)

// HandlerWrapper TODO
type HandlerWrapper struct {
	App *ctx.ApplicationContext
}

// Wrap TODO
func (wp *HandlerWrapper) Wrap(f CtrlFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		f(c, wp.App)
	}
}
