package main

import (
	"github.com/Ghamster0/os-rq-fsender/app/ctx"
	"github.com/Ghamster0/os-rq-fsender/app/router"
	"github.com/gin-gonic/gin"
)

func main() {
	app, _ := ctx.StartApplication()
	r := gin.Default()
	router.InitRouter(r, app)
	r.Run("0.0.0.0:8080")
}
