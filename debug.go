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
	items := []Item{
		{Name: "banana", Price: 80, Comment: "delicious and nutritious", Genre_id: 1},
		{Name: "orange", Price: 120, Comment: "juicy and refreshing", Genre_id: 1},
		{Name: "grape", Price: 200, Comment: "sweet and flavorful", Genre_id: 1},
		{Name: "kiwi", Price: 150, Comment: "exotic and tangy", Genre_id: 1},
		{Name: "pineapple", Price: 300, Comment: "tropical and juicy", Genre_id: 1},
		{Name: "watermelon", Price: 500, Comment: "refreshing and hydrating", Genre_id: 1},
		{Name: "peach", Price: 180, Comment: "aromatic and juicy", Genre_id: 1},
		{Name: "pear", Price: 120, Comment: "crisp and juicy", Genre_id: 1},
		{Name: "mango", Price: 250, Comment: "juicy and tropical", Genre_id: 1},
		{Name: "human", Price: 0, Comment: "inedible and unacceptable", Genre_id: 6},
		{Name: "carrot", Price: 90, Comment: "crunchy and nutritious", Genre_id: 2},
		{Name: "broccoli", Price: 130, Comment: "nutritious and delicious", Genre_id: 2},
		{Name: "celery", Price: 100, Comment: "crisp and refreshing", Genre_id: 2},
		{Name: "tomato", Price: 70, Comment: "juicy and flavorful", Genre_id: 2},
		{Name: "cucumber", Price: 80, Comment: "refreshing and crunchy", Genre_id: 2},
		{Name: "lemonade", Price: 150, Comment: "refreshing and tangy", Genre_id: 3},
		{Name: "iced tea", Price: 120, Comment: "cool and refreshing", Genre_id: 3},
		{Name: "orange juice", Price: 200, Comment: "fresh and flavorful", Genre_id: 3},
		{Name: "potato chips", Price: 90, Comment: "crispy and salty", Genre_id: 4},
		{Name: "popcorn", Price: 100, Comment: "crunchy and buttery", Genre_id: 4},
		{Name: "cheese", Price: 180, Comment: "creamy and flavorful", Genre_id: 5},
		{Name: "yogurt", Price: 120, Comment: "creamy and tangy", Genre_id: 5},
		{Name: "beef", Price: 300, Comment: "juicy and flavorful", Genre_id: 6},
		{Name: "chicken", Price: 250, Comment: "tender and juicy", Genre_id: 6},
		{Name: "cherry", Price: 170, Comment: "juicy and sweet", Genre_id: 1},
		{Name: "strawberry", Price: 120, Comment: "sweet and juicy", Genre_id: 1},
		{Name: "potato", Price: 60, Comment: "versatile and nutritious", Genre_id: 2},
		{Name: "lettuce", Price: 70, Comment: "crisp and fresh", Genre_id: 2},
		{Name: "spinach", Price: 100, Comment: "nutritious and versatile", Genre_id: 2},
		{Name: "iced coffee", Price: 180, Comment: "refreshing and energizing", Genre_id: 3},
		{Name: "soda", Price: 100, Comment: "carbonated and sweet", Genre_id: 3},
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
