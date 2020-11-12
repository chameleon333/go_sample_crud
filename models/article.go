package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Article is struct
type Article struct {
	ID        int
	Title     string
	Body      []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

var db gorm.DB

func LoadArticle(id int) (*Article, error) {
	var article Article

	db.First(&article, "id=?", id)

	title := article.Title
	body := article.Body

	return &Article{ID: id, Title: title, Body: body}, nil
}

func LoadAllArticle() []Article {
	var articles []Article

	db.Order("created_at desc").Find(&articles)

	return articles
}

func (article *Article) Save() {
	title := article.Title
	body := article.Body

	db.First(&article, "id=?", article.ID)
	article.Title = title
	article.Body = body

	db.Save(&article)
}

func (article *Article) Delete() {
	db.First(&article, "id=?", article.ID)
	db.Delete(&article)
}
