package main

import (
	"log"
	"work/controllers"
	"work/db"
	"work/models"
	"work/session"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	db.GormConnect()
	db.DB.AutoMigrate(&models.User{}, &models.Article{})
}

func main() {

	router := gin.Default()

	// セッションを設定
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	router.LoadHTMLGlob("view/*/**")

	auth := router.Group("/")
	auth.Use(session.SessionCheck())
	{
		auth.GET("/new", controllers.NewArticles)
		auth.GET("/edit/:id", controllers.EditArticles)
		auth.GET("/logout", controllers.Logout)
		auth.GET("/delete/:id", controllers.DeleteArticle)
		auth.POST("/save", controllers.SaveArticles)
	}

	router.GET("/list", controllers.ListArticles)
	router.GET("/view/:id", controllers.ViewArticles)
	router.GET("/login", controllers.LoginFormRoute)
	router.POST("/login", controllers.Login)
	router.GET("/signup", controllers.SignUpFormRoute)
	router.POST("/signup", controllers.SignUp)
	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err.Error())
	}
}
