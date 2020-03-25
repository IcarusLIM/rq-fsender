package main

import (
	"github.com/Ghamster0/os-rq-fsender/app/ctx"
	"github.com/Ghamster0/os-rq-fsender/app/router"
	"github.com/gin-gonic/gin"
)

func main() {
	ctx := &ctx.ApplicationContext{}
	r := gin.Default()
	router.InitRouter(r, ctx)
	r.Run()
}
