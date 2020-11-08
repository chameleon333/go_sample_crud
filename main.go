package main

import (
	"work/controllers"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
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

func main() {

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/new", controllers.NewArticles)
	router.GET("/list", controllers.ListArticles)
	router.GET("/view/:id", controllers.ViewArticles)
	router.GET("/edit/:id", controllers.EditArticles)
	router.GET("/delete/:id", controllers.DeleteArticle)
	router.POST("/save", controllers.SaveArticles)
	router.GET("/signin", SigninFormRoute)
	router.POST("/signin", Signin)
	router.Run(":8080")
}
