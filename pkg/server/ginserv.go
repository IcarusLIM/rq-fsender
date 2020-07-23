package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func NewEngine(conf *viper.Viper) *gin.Engine {
	debug := false
	if conf.IsSet("debug") {
		debug = conf.GetBool("debug")
	}
	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	return engine
}

func EnableCROS(r *gin.RouterGroup) {
	r.Use(
		func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Headers", "*")
			c.Header("Access-Control-Allow-Methods", "*")
			if c.Request.Method == http.MethodOptions {
				c.Status(http.StatusNoContent)
			} else {
				c.Next()
			}
		},
	)
}

func NewRouterGroup(conf *viper.Viper, engine *gin.Engine) *gin.RouterGroup {
	return engine.Group(conf.GetString("http.api.path"))
}
