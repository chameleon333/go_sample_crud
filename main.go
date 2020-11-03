package main

import (
	"html/template"
	"log"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/kelseyhightower/envconfig"
)

// Article is struct
type Article struct {
	ID    int
	Title string
	Body  []byte
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
	db, err := gorm.Open(goenv.DB_CONNECTION, CONNECT)
	if err != nil {
		panic(err)
	}
	return db
}

func (article *Article) save() {
	db := gormConnect()
	defer db.Close()
	body := article.Body

	db.First(&article, "title=?", article.Title)
	article.Body = body

	db.Save(&article)
}

func (article *Article) delete() {
	db := gormConnect()
	defer db.Close()

	db.First(&article, "title=?", article.Title)
	db.Delete(&article)
}

func loadArticle(title string) (*Article, error) {
	var article Article

	db := gormConnect()
	defer db.Close()

	db.First(&article, "title=?", title)

	body := article.Body
	id := article.ID

	return &Article{ID: id, Title: title, Body: body}, nil
}

func loadAllArticle() []Article {
	var articles []Article

	db := gormConnect()
	defer db.Close()

	db.Order("created_at desc").Find(&articles)

	return articles
}

func renderTemplate(writer http.ResponseWriter, tmpl string, article interface{}) {
	t, _ := template.ParseFiles(tmpl + ".html")
	err := t.Execute(writer, article)
	if err != nil {
		panic(err.Error())
	}
}

func listHandler(writer http.ResponseWriter, r *http.Request) {
	articles := loadAllArticle()
	renderTemplate(writer, "list", articles)
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	article, err := loadArticle(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", article)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	article, err := loadArticle(title)
	if err != nil {
		article = &Article{Title: title}
	}
	renderTemplate(w, "edit", article)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	article := &Article{Title: title, Body: []byte(body)}
	article.save()

	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func deleteHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	article := &Article{Title: title, Body: []byte(body)}
	article.delete()
	http.Redirect(w, r, "/list/", http.StatusFound)
}

var validPath = regexp.MustCompile("^/(list|edit|save|view|delete)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func ListArticles(c *gin.Context) {
	articles := loadAllArticle()
	c.HTML(200, "list.html", gin.H{
		"articles": articles,
	})
}

func EditArticles(c *gin.Context) {
	title := c.Param("title")
	article, err := loadArticle(title)
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
	title := c.PostForm("title")
	body := c.PostForm("body")
	article := &Article{Title: title, Body: []byte(body)}
	article.save()
	c.Redirect(http.StatusSeeOther, "/list")
}

func DeleteArticle(c *gin.Context) {
	title := c.Param("title")
	db := gormConnect()
	defer db.Close()
	article := &Article{Title: title}
	article.delete()
	c.Redirect(http.StatusSeeOther, "/list")
}

func ViewArticles(c *gin.Context) {
	title := c.Param("title")
	article, err := loadArticle(title)
	if err != nil {
		// リダイレクトする
		c.Redirect(http.StatusNotFound, "404.html")
		return
	}
	c.HTML(200, "view.html", gin.H{
		"article": article,
	})
}

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello world",
		})
	})
	router.LoadHTMLGlob("templates/*")
	router.GET("/list", ListArticles)
	router.GET("/view/:title", ViewArticles)
	router.GET("/edit/:title", EditArticles)
	router.GET("/delete/:title", DeleteArticle)
	router.POST("/save", SaveArticles)

	router.Run(":8080")

	// http.HandleFunc("/list/", listHandler)
	// http.HandleFunc("/view/", makeHandler(viewHandler))
	// http.HandleFunc("/edit/", makeHandler(editHandler))
	// http.HandleFunc("/save/", makeHandler(saveHandler))
	// http.HandleFunc("/delete/", makeHandler(deleteHandler))
	// log.Fatal(http.ListenAndServe(":8080", nil))
}
