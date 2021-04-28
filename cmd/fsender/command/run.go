package command

import (
	"context"
	"net/http"
	"os"

	"github.com/Ghamster0/os-rq-fsender/app/controller"
	"github.com/Ghamster0/os-rq-fsender/app/router"
	"github.com/Ghamster0/os-rq-fsender/pkg/command"
	"github.com/Ghamster0/os-rq-fsender/pkg/config"
	"github.com/Ghamster0/os-rq-fsender/pkg/db"
	"github.com/Ghamster0/os-rq-fsender/pkg/global"
	logconf "github.com/Ghamster0/os-rq-fsender/pkg/log"
	"github.com/Ghamster0/os-rq-fsender/pkg/server"
	"github.com/Ghamster0/os-rq-fsender/send/task"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func init() {
	Root.AddCommand(command.NewRunCommand("fsender", run))
}

func run(conf *viper.Viper) {
	newConfig := func() (*viper.Viper, error) {
		err := config.LoadConfig(conf, global.EnvPrefix, global.DefaultConfig)
		return conf, err
	}

	runServer := func(lifecycle fx.Lifecycle, conf *viper.Viper, engine *gin.Engine) {
		lifecycle.Append(
			fx.Hook{
				OnStart: func(ctx context.Context) error {
					go http.ListenAndServe("0.0.0.0:"+conf.GetString("http.port"), engine)
					return nil
				},
			},
		)
	}

	initDirs := func(lifecycle fx.Lifecycle, conf *viper.Viper) {
		ensureDir := func(p string) {
			if _, err := os.Stat(p); os.IsNotExist(err) {
				os.Mkdir(p, os.ModePerm)
			}
		}
		lifecycle.Append(
			fx.Hook{
				OnStart: func(ctx context.Context) error {
					ensureDir(conf.GetString("upload.file"))
					ensureDir(conf.GetString("upload.log"))
					return nil
				},
			},
		)
	}

	app := fx.New(
		fx.Provide(
			newConfig,
			db.NewDB,
			server.NewEngine,
			server.NewRouterGroup,
			controller.NewSendController,
			controller.NewBatchController,
			task.NewTaskBox,
			task.NewFileService,
			task.NewBatchService,
			task.NewCleanService,
		),
		fx.Invoke(
			initDirs,
			task.CleanServiceServ,
			logconf.ConfigLogger,
			task.TaskBoxServe,
			server.EnableCROS,
			router.FileRouter,
			router.BatchRouter,
			runServer,
		),
	)
	app.Run()
}
