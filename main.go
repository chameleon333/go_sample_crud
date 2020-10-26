package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
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

// func gormConnect() *gorm.DB {
// 	var goenv Env
// 	err := envconfig.Process("", &goenv)
// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}
// 	PROTOCOL := "tcp(" + goenv.DB_HOST + ":" + goenv.DB_PORT + ")"
// 	CONNECT := goenv.DB_USERNAME + ":" + goenv.DB_PASSWORD + "@" + PROTOCOL + "/" + goenv.DB_DATABASE
// 	db, err := gorm.Open(goenv.DB_CONNECTION, CONNECT)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return db
// }

func DBConnect() (db *sql.DB) {
	var goenv Env
	err2 := envconfig.Process("", &goenv)
	if err2 != nil {
		panic(err2.Error())
	}
	PROTOCOL := "tcp(" + goenv.DB_HOST + ":" + goenv.DB_PORT + ")"

	CONNECT := goenv.DB_USERNAME + ":" + goenv.DB_PASSWORD + "@" + PROTOCOL + "/" + goenv.DB_DATABASE
	db, err := sql.Open("mysql", CONNECT)
	if err != nil {
		panic(err)
	}
	return db
}

func (p *Article) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func (p *Article) delete() error {
	filename := p.Title + ".txt"
	fmt.Println(filename)
	return os.Remove(filename)
}

func loadArticle(title string) (*Article, error) {
	// db := DBConnect()
	// defer db.Close()
	// var body string
	// // stmtOut, err := db.Prepare(fmt.Sprintf("SELECT * FROM articles where title = ?", title)) //
	// err := db.QueryRow("SELECT * FROM articles where title = ?", title).Scan(&body)
	// if err != nil {
	// 	panic(err.Error())
	// }
	// fmt.Println(body)
	// fmt.Println(title)

	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Article{Title: title, Body: body}, nil
}

func renderTemplate(writer http.ResponseWriter, tmpl string, article *Article) {
	t, _ := template.ParseFiles(tmpl + ".html")
	err := t.Execute(writer, article)
	if err != nil {
		panic(err.Error())
	}
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
	err := article.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func deleteHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	article := &Article{Title: title, Body: []byte(body)}
	err := article.delete()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

var validPath = regexp.MustCompile("^/(edit|save|view|delete)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		fmt.Println(m)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func main() {

	db := DBConnect()
	// db := gormConnect()
	defer db.Close()
	var article Article
	// var id string
	// var title string
	// var body string

	article.ID = 1
	// db.First(&article)
	// stmtOut, err := db.Prepare(fmt.Sprintf("SELECT * FROM articles where title = ?", title)) //
	err := db.QueryRow("SELECT * FROM articles where title = ?", "title1").Scan(&article.ID, &article.Title, &article.Body)
	// err := db.QueryRow("SELECT * FROM articles where title = ?", "title1")
	// fmt.Println(err)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(article.ID)
	fmt.Println(article.Title)
	fmt.Println(article.Body)
	fmt.Println("---------")
	// fmt.Println(title)

	// db := DBConnect()
	// defer db.Close()

	// // rows, err := db.Query("show tables") //
	// rows, err := db.Query("SELECT * FROM articles") //
	// if err != nil {
	// 	panic(err.Error())
	// }
	// fmt.Println(rows)

	// columns, err := rows.Columns() // カラム名を取得
	// if err != nil {
	// 	panic(err.Error())
	// }

	// fmt.Println(columns)

	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/delete/", makeHandler(deleteHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
