package controllers

import (
	"net/http"
	"strconv"
	"work/models"

	"github.com/gin-gonic/gin"
)

func TopArticles(c *gin.Context) {
	c.Redirect(http.StatusSeeOther, "/list")
}

func NewArticles(c *gin.Context) {
	c.HTML(200, "edit.html", gin.H{})
}

func ListArticles(c *gin.Context) {
	articles := models.LoadAllArticle()
	c.HTML(200, "list.html", gin.H{
		"articles": articles,
	})
}

func EditArticles(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	article, err := models.LoadArticle(id)
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
	article := &models.Article{ID: id, Title: title, Body: []byte(body)}
	article.Save()
	c.Redirect(http.StatusSeeOther, "/list")
}

func DeleteArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	article := &models.Article{ID: id}
	article.Delete()
	c.Redirect(http.StatusSeeOther, "/list")
}

func ViewArticles(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	article, err := models.LoadArticle(id)
	if err != nil {
		// リダイレクトする
		c.Redirect(http.StatusNotFound, "404.html")
		return
	}
	c.HTML(200, "view.html", gin.H{
		"article": article,
	})
}
