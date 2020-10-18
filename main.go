package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kelseyhightower/envconfig"
)

type Page struct {
	Title string
	Body  []byte
}

type Env struct {
	DB_PORT        int
	MYSQL_DATABASE string
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func (p *Page) delete() error {
	filename := p.Title + ".txt"
	fmt.Println(filename)
	return os.Remove(filename)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles(tmpl + ".html")
	t.Execute(w, p)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(len("/delete/"))
	title := r.URL.Path[len("/delete/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.delete()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
	var goenv Env
	envconfig.Process("", &goenv)
	fmt.Println(goenv.DB_PORT)
	fmt.Println(goenv.MYSQL_DATABASE)

	db, err := sql.Open("mysql", "admin:secret@tcp(db-host:3306)/test")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// rows, err := db.Query("show tables") //
	rows, err := db.Query("SELECT * FROM articles") //
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(rows)

	columns, err := rows.Columns() // カラム名を取得
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(columns)

	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/delete/", deleteHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
