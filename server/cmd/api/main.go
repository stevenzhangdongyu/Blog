package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"personal-blog/server/internal/article"
)

func main() {
	db, err := article.GetDBConnection()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	router := gin.Default()

	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/api/articles", func(c *gin.Context) {
		articles, err := article.GetArticles(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		//序列化
		c.JSON(http.StatusOK, gin.H{"articles": articles})
	})

	router.GET("/api/articles/:slug", func(c *gin.Context) {
		slug := c.Param("slug")
		foundArticle, err := article.GetArticleBySlug(db, slug)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if foundArticle == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
			return
		}
		c.JSON(http.StatusOK, foundArticle)
	})

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
