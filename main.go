package main

import (
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Item struct {
	Name  string `gorm:"primaryKey"`
	Price int    `gorm:"not null"`
}

func main() {
	db, err := gorm.Open(sqlite.Open("items.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.Table("items").AutoMigrate(&Item{})

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		rawHeaders, _ := httputil.DumpRequest(c.Request, false)
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Rawheader": string(rawHeaders),
		})
	})
	r.GET("/items", func(c *gin.Context) {
		var items []Item
		db.Find(&items)

		c.JSON(http.StatusOK, gin.H{"items": items})
	})
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	r.POST("/", func(c *gin.Context) {
		var item Item
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := db.Create(&item).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, item)
	})

	r.Run(":8080")
}
