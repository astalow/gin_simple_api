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
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"not null"`
	Price    uint   `gorm:"not null"`
	Genre_id uint   `gorm:"not null"`
	Comment  string
}

type Genre struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"not null"`
}

func migrateDB() *gorm.DB {
	// DBインスタンスを初期化する
	db, err := gorm.Open(sqlite.Open("items.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// itemsテーブルをマイグレーションする
	if db.Migrator().HasTable(&Item{}) && db.Migrator().HasTable(&Genre{}) {
		db.Exec("DROP TABLE items")
		db.Exec("DROP TABLE genres")
	}
	if db.AutoMigrate(&Item{}) != nil || db.AutoMigrate(&Genre{}) != nil {
		panic("failed to migrate database")
	}
	setDebugItems(db)
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
			<th>Comment</th>
		</tr>
		{{range .}}
			<tr>
				<td>{{.Name}}</td>
				<td>{{.Price}}</td>
				<td>{{.Comment}}</td>
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
