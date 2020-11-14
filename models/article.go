package models

import (
	"time"
	"work/db"
)

// Article is struct
type Article struct {
	ID        int
	Title     string
	Body      []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

func LoadArticle(id int) (*Article, error) {
	var article Article

	db.DB.First(&article, "id=?", id)

	title := article.Title
	body := article.Body

	return &Article{ID: id, Title: title, Body: body}, nil
}

func LoadAllArticle() []Article {
	var articles []Article

	db.DB.Order("created_at desc").Find(&articles)

	return articles
}

func (article *Article) Save() {
	title := article.Title
	body := article.Body

	db.DB.First(&article, "id=?", article.ID)
	article.Title = title
	article.Body = body

	db.DB.Save(&article)
}

func (article *Article) Delete() {
	db.DB.First(&article, "id=?", article.ID)
	db.DB.Delete(&article)
}
