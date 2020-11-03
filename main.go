package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/kelseyhightower/envconfig"
)

// Article is struct
type Article struct {
	ID        int
	Title     string
	Body      []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Env is struct
type Env struct {
	DB_CONNECTION string
	DB_PORT       string
	DB_HOST       string
	DB_USERNAME   string
	DB_PASSWORD   string
	DB_DATABASE   string
}

func gormConnect() *gorm.DB {
	var goenv Env
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

func (article *Article) save() {
	db := gormConnect()
	defer db.Close()
	title := article.Title
	body := article.Body

	db.First(&article, "id=?", article.ID)
	article.Title = title
	article.Body = body

	db.Save(&article)
}

func (article *Article) delete() {
	db := gormConnect()
	defer db.Close()

	db.First(&article, "id=?", article.ID)
	db.Delete(&article)
}

func loadArticle(id int) (*Article, error) {
	var article Article

	db := gormConnect()
	defer db.Close()

	db.First(&article, "id=?", id)

	title := article.Title
	body := article.Body

	return &Article{ID: id, Title: title, Body: body}, nil
}

func loadAllArticle() []Article {
	var articles []Article

	db := gormConnect()
	defer db.Close()

	db.Order("created_at desc").Find(&articles)

	return articles
}

func ListArticles(c *gin.Context) {
	articles := loadAllArticle()
	c.HTML(200, "list.html", gin.H{
		"articles": articles,
	})
}

func EditArticles(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	article, err := loadArticle(id)
	if err != nil {
		// リダイレクトする
		c.Redirect(http.StatusNotFound, "404.html")
		return
	}
	c.HTML(200, "edit.html", gin.H{
		"article": article,
	})
}

func SaveArticles(c *gin.Context) {
	id, _ := strconv.Atoi(c.PostForm("id"))
	title := c.PostForm("title")
	body := c.PostForm("body")
	article := &Article{ID: id, Title: title, Body: []byte(body)}
	article.save()
	c.Redirect(http.StatusSeeOther, "/list")
}

func DeleteArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	db := gormConnect()
	defer db.Close()
	article := &Article{ID: id}
	article.delete()
	c.Redirect(http.StatusSeeOther, "/list")
}

func ViewArticles(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	article, err := loadArticle(id)
	if err != nil {
		// リダイレクトする
		c.Redirect(http.StatusNotFound, "404.html")
		return
	}
	c.HTML(200, "view.html", gin.H{
		"article": article,
	})
}

func TopArticles(c *gin.Context) {
	c.Redirect(http.StatusSeeOther, "/list")
}

func NewArticles(c *gin.Context) {
	c.HTML(200, "edit.html", gin.H{})
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/", TopArticles)
	router.GET("/new", NewArticles)
	router.GET("/list", ListArticles)
	router.GET("/view/:id", ViewArticles)
	router.GET("/edit/:id", EditArticles)
	router.GET("/delete/:id", DeleteArticle)
	router.POST("/save", SaveArticles)
	router.Run(":8080")
}
