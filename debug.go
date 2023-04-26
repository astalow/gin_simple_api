package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func getUserAndPassword(c *gin.Context) (User, error) {
	var user User

	if err := c.ShouldBindJSON(&user); err != nil {
		fmt.Printf("JSON parsing error occurred: %v\n", err)
		return user, err
	}
	return user, nil
}

func handlePasswordComparison(r *gin.Engine) {
	r.GET("/compare-password", func(c *gin.Context) {
		c.HTML(http.StatusOK, "compare-password.html", nil)
	})
	r.POST("/compare-password", func(c *gin.Context) {
		user, err := getUserAndPassword(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("password: ", user.Password)
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return
		}
		fmt.Println("Hashed password: ", string(hashedPassword))

		err = bcrypt.CompareHashAndPassword(hashedPassword, []byte("password"))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
			return
		}

		fmt.Println("Password is valid")
	})
}

func setDebugItems(db *gorm.DB) {
	com := "oisii"
	items := []Item{
		{Name: "banana", Price: 80, Comment: com, Genre_id: 1},
		{Name: "orange", Price: 120, Comment: com, Genre_id: 1},
		{Name: "grape", Price: 200, Comment: com, Genre_id: 1},
		{Name: "kiwi", Price: 150, Comment: com, Genre_id: 1},
		{Name: "pineapple", Price: 300, Comment: com, Genre_id: 1},
		{Name: "watermelon", Price: 500, Comment: com, Genre_id: 1},
		{Name: "peach", Price: 180, Comment: com, Genre_id: 1},
		{Name: "pear", Price: 120, Comment: com, Genre_id: 1},
		{Name: "mango", Price: 250, Comment: com, Genre_id: 1},
		{Name: "human", Price: 0, Comment: "oisikunai", Genre_id: 1},
	}

	for _, item := range items {
		db.Create(&item)
	}
}

func setDebugGenres(db *gorm.DB) {
	items := []Genre{
		{
			Name: "fruit",
		},
		{
			Name: "vegetable",
		},
		{
			Name: "drink",
		},
		{
			Name: "snack",
		},
		{
			Name: "dairy",
		},
		{
			Name: "meat",
		},
	}

	for _, item := range items {
		db.Create(&item)
	}
}
