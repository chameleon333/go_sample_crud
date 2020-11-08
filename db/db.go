package db

import (
	"log"
	"work/env"

	"github.com/jinzhu/gorm"
	"github.com/kelseyhightower/envconfig"
)

func GormConnect() *gorm.DB {
	var goenv env.Env
	err := envconfig.Process("", &goenv)
	if err != nil {
		log.Fatal(err.Error())
	}
	PROTOCOL := "tcp(" + goenv.DB_HOST + ":" + goenv.DB_PORT + ")"
	CONNECT := goenv.DB_USERNAME + ":" + goenv.DB_PASSWORD + "@" + PROTOCOL + "/" + goenv.DB_DATABASE
	OPTION := "?parseTime=true&loc=Asia%2FTokyo"
	db, err := gorm.Open(goenv.DB_CONNECTION, CONNECT+OPTION)
	if err != nil {
		panic(err)
	}
	return db
}
