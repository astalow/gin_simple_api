package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Account struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"not null"`
	Password string `gorm:"not null"`
}

func migrateAccountDB() *gorm.DB {
	// DBインスタンスを初期化する
	db, err := gorm.Open(sqlite.Open("account.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// テーブルをマイグレーションする
	db.Exec("DELETE FROM accounts")
	err = db.AutoMigrate(&Account{})
	if err != nil {
		panic("failed to migrate database")
	}

	accounts := []Account{
		{ID: 1, Username: "user", Password: "password"},
	}

	for _, account := range accounts {
		db.Create(&account)
	}
	return db
}
