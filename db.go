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
	setDebugGenres(db)
	return db
}

func buildHTMLTableFromDB(db *gorm.DB) string {
	var items []struct {
		ID         uint
		Name       string
		Price      uint
		Genre_name string
		Comment    string
	}

	if err := db.Raw("SELECT i.name, i.price, i.comment,COALESCE(g.name, 'other') AS genre_name	FROM items i	LEFT JOIN genres g ON i.genre_id = g.id	").Scan(&items).Error; err != nil {
		return ""
	}

	const tableTemplate = `
	<table>
		<tr>
			<th>Name</th>
			<th>Price</th>
			<th>Genre</th>
			<th>Comment</th>
		</tr>
		{{range .}}
			<tr>
				<td>{{.Name}}</td>
				<td>{{.Price}}</td>
				<td>{{.Genre_name}}</td>
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

func buildHTMLTableFromDB2(db *gorm.DB) string {
	var avg []struct {
		Name         string
		AveragePrice string
	}

	q := `SELECT genres.name, AVG(items.price) AS AveragePrice
		FROM genres
		LEFT JOIN items ON genres.id = items.genre_id
		GROUP BY genres.name;`

	if err := db.Raw(q).Scan(&avg).Error; err != nil {
		return ""
	}

	const tableTemplate = `
	<table>
		<tr>
			<th>Genre</th>
			<th>AveragePrice</th>
		</tr>
		{{range .}}
			<tr>
				<td>{{.Name}}</td>
				<td>{{.AveragePrice}}</td>
			</tr>
		{{end}}
	</table>`

	tmpl, err := template.New("table2").Parse(tableTemplate)
	if err != nil {
		return ""
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, avg); err != nil {
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
