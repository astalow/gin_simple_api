package main

import (
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	Name  string `gorm:"primaryKey"`
	Price int    `gorm:"not null"`
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&User{})

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		rawHeaders, _ := httputil.DumpRequest(c.Request, false)
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Rawheader": string(rawHeaders),
		})
	})

	r.POST("/", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := db.Create(&user).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, user)
	})

	r.Run(":8080")
}
