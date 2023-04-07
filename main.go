package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	r.POST("/submit", func(c *gin.Context) {
		name := c.PostForm("name")
		price := c.PostForm("price")
		c.JSON(http.StatusOK, gin.H{
			"name":  name,
			"price": price,
		})
	})

	r.Run()
}
