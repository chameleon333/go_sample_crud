package main

import (
	"log"
	"work/controllers"
	"work/env"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/kelseyhightower/envconfig"
)

func SigninFormRoute(c *gin.Context) {
	c.HTML(200, "signin.html", gin.H{})
}

func Signin(c *gin.Context) {
	c.HTML(200, "signin.html", gin.H{})

	// isExist, user := HTTPRequestManager.IsLoginUserExist(c)

	// if isExist {
	// 	SessionManager.Login(c, user)
	// }
	// info := SessionManager.GetSessionInfo(c)

	// server.SetHTMLTemplate(templates["index"])
	// c.HTML(200, "_base.html", gin.H{
	// 	"SessionInfo": info,
	// })
}

func init() {
	var goenv env.Env
	err := envconfig.Process("", &goenv)
	if err != nil {
		log.Fatal(err.Error())
	}
	PROTOCOL := "tcp(" + goenv.DB_HOST + ":" + goenv.DB_PORT + ")"
	CONNECT := goenv.DB_USERNAME + ":" + goenv.DB_PASSWORD + "@" + PROTOCOL + "/" + goenv.DB_DATABASE
	OPTION := "?parseTime=true&loc=Asia%2FTokyo"
	db, err = gorm.Open(goenv.DB_CONNECTION, CONNECT+OPTION)
	if err != nil {
		panic(err)
	}
}

var db *gorm.DB

func main() {

	router := gin.Default()
	router.LoadHTMLGlob("view/*/**")
	router.GET("/new", controllers.NewArticles)
	router.GET("/list", controllers.ListArticles)
	router.GET("/view/:id", controllers.ViewArticles)
	router.GET("/edit/:id", controllers.EditArticles)
	router.GET("/delete/:id", controllers.DeleteArticle)
	router.POST("/save", controllers.SaveArticles)
	router.GET("/signin", SigninFormRoute)
	router.POST("/signin", Signin)
	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err.Error())
	}
}
