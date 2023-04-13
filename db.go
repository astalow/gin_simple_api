package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
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

func buildHTMLTableFromDB(db *gorm.DB) string {
	var items []Item
	if err := db.Find(&items).Error; err != nil {
		return ""
	}

	const tableTemplate = `
	<table>
		<tr>
			<th>Name</th>
			<th>Price</th>
		</tr>
		{{range .}}
			<tr>
				<td>{{.Name}}</td>
				<td>{{.Price}}</td>
			</tr>
		{{end}}
	</table>`

	tmpl, err := template.New("table").Parse(tableTemplate)
	if err != nil {
		return ""
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, items); err != nil {
		return ""
	}

	return buf.String()
}

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// リクエストボディをログに出力
		headers := c.Request.Header
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
			return
		}
		// リクエストbodyを元に戻す
		defer func() {
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}()

		// bodyをそのまま出力する
		requestInfo := fmt.Sprintf("[Request info]\nMethod: %s\nURL: %s\nHeaders: %v\nBody: %s\n\n", c.Request.Method, c.Request.URL, headers, string(body))
		fmt.Print(requestInfo)

		c.Next()
	}
}

func printDB(db *gorm.DB) {
	var users []Item
	db.Find(&users)

	fmt.Println()
	for _, user := range users {
		fmt.Printf("name: %s, price: %d\n", user.Name, user.Price)
	}
	fmt.Println()
}