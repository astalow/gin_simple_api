package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main2() {
	// SQLiteデータベースに接続
	db, err := sql.Open("sqlite3", "sample.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// itemsテーブルを作成する
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		price INTEGER
	);`)
	if err != nil {
		log.Fatal(err)
	}

	// jenresテーブルを作成する
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS jenres (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT
	);`)
	if err != nil {
		log.Fatal(err)
	}

	// ginルーターを設定する
	r := gin.Default()

	// ルートパスにアクセスした時の処理を設定する
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World")
	})

	// サーバーを起動する
	err = r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
