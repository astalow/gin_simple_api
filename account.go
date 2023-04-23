package main

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Account struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"not null"`
	Password []byte `gorm:"not null"`
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	if err != nil {
		return nil
	}

	accounts := []Account{
		{ID: 1, Username: "kudamonodaisuki", Password: hashedPassword},
	}

	for _, account := range accounts {
		db.Create(&account)
	}
	return db
}
