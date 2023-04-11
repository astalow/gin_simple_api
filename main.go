package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Item struct {
	// ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"not null"`
	Price uint   `gorm:"not null"`
}

func migrateDB() *gorm.DB {
	// DBインスタンスを初期化する
	db, err := gorm.Open(sqlite.Open("items.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// itemsテーブルをマイグレーションする
	db.Exec("DELETE FROM items")
	err = db.AutoMigrate(&Item{})
	if err != nil {
		panic("failed to migrate database")
	}

	items := []Item{
		{Name: "banana", Price: 80},
		{Name: "orange", Price: 120},
		{Name: "grape", Price: 200},
		{Name: "kiwi", Price: 150},
		{Name: "pineapple", Price: 300},
		{Name: "watermelon", Price: 500},
		{Name: "peach", Price: 180},
		{Name: "pear", Price: 120},
		{Name: "mango", Price: 250},
		{Name: "human", Price: 0},
	}

	for _, item := range items {
		db.Create(&item)
	}
	return db
}

func printBody(c *gin.Context) {
	// リクエストボディをログに出力
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
		return
	}
	// リクエストbodyを元に戻す
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	// bodyをそのまま出力する
	fmt.Print("\n" + string(body) + "\n\n")
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	db := migrateDB()

	r.GET("/", func(c *gin.Context) {
		var items []Item
		db.Find(&items)

		c.JSON(http.StatusOK, gin.H{"items": items})
	})

	r.GET("/view", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.POST("/", func(c *gin.Context) {
		printBody(c)
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
