package db

import (
	"log"
	"work/env"

	"github.com/jinzhu/gorm"
	"github.com/kelseyhightower/envconfig"
)

var DB *gorm.DB

func GormConnect() {
	var goenv env.Env
	err := envconfig.Process("", &goenv)
	if err != nil {
		log.Fatal(err.Error())
	}
	PROTOCOL := "tcp(" + goenv.DB_HOST + ":" + goenv.DB_PORT + ")"
	CONNECT := goenv.DB_USERNAME + ":" + goenv.DB_PASSWORD + "@" + PROTOCOL + "/" + goenv.DB_DATABASE
	OPTION := "?parseTime=true&loc=Asia%2FTokyo"
	DB, err = gorm.Open(goenv.DB_CONNECTION, CONNECT+OPTION)
	if err != nil {
		panic(err)
	}
}
