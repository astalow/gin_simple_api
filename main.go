package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	db := migrateDB()
	printDB(db)

	handleUserSession(r)
	// handlePasswordComparison(r)

	r.GET("/", func(c *gin.Context) {
		var items []Item
		db.Find(&items)

		c.JSON(http.StatusOK, gin.H{"items": items})
	})

	r.GET("/view", func(c *gin.Context) {
		// テンプレートファイルをパースする
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			log.Fatal(err)
		}

		// データベースからHTMLテーブルを構築する
		htmlTable := buildHTMLTableFromDB(db)

		// テンプレートをレンダリングする
		err = tmpl.Execute(c.Writer, gin.H{
			"table": template.HTML(htmlTable),
		})
		if err != nil {
			log.Fatal(err)
		}
	})

	r.POST("/", func(c *gin.Context) {
		// RequestLogger(c)
		var item Item
		if err := c.ShouldBindJSON(&item); err != nil {
			// JSONパースエラーが発生した場合
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		if err := db.Create(&item).Error; err != nil {
			// データベースエラーが発生した場合
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item"})
			return
		}

		c.JSON(http.StatusOK, item)
	})
	r.Run(":8080")
}
