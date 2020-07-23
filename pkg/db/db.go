package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
)

func NewDB(conf *viper.Viper) *gorm.DB {
	host := conf.GetString("db.mysql.host")
	user := conf.GetString("db.mysql.user")
	password := conf.GetString("db.mysql.password")
	dbname := conf.GetString("db.mysql.dbname")
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@%s/%s?charset=utf8&parseTime=True&loc=Local&timeout=10s", user, password, host, dbname))
	if err != nil {
		panic("Can't connect to db: " + fmt.Sprintf("%s:%s@%s/%s", user, password, host, dbname) + "\n" + err.Error())
	}
	return db
}
