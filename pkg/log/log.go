package logconf

import (
	"log"
	"os"

	"github.com/op/go-logging"
	"github.com/spf13/viper"
)

func ConfigLogger(conf *viper.Viper) {
	logging.SetBackend(logging.NewLogBackend(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile))
	level := "debug"
	if conf.IsSet("log.level") {
		level = conf.GetString("log.level")
	}
	logLevel, err := logging.LogLevel(level)
	if err != nil {
		logLevel = logging.DEBUG
	}
	logging.SetLevel(logLevel, "")
}
