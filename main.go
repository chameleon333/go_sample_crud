package main

import (
	"log"
	"work/controllers"
	"work/db"
	"work/models"

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
	router.GET("/new", controllers.NewArticles)
	router.GET("/list", controllers.ListArticles)
	router.GET("/view/:id", controllers.ViewArticles)
	router.GET("/edit/:id", controllers.EditArticles)
	router.GET("/delete/:id", controllers.DeleteArticle)
	router.POST("/save", controllers.SaveArticles)
	router.GET("/signin", controllers.SignInFormRoute)
	router.POST("/signin", controllers.SignIn)
	router.GET("/signup", controllers.SignUpFormRoute)
	router.POST("/signup", controllers.SignUp)
	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err.Error())
	}
}
