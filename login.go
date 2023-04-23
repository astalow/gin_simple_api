package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string
	Password string
}

func handleUserSession(r *gin.Engine) {
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))
	db := migrateAccountDB()

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	r.POST("/login", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			// JSONパースエラーが発生した場合
			fmt.Printf("JSON parsing error occurred: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		// フォームから送信されたユーザー名とパスワードを取得する。
		username := user.Username
		password := user.Password

		var ac []Account

		err := db.Where("username LIKE ?", username).Find(&ac).Error
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to query users"})
			return
		}

		if len(ac) == 0 {
			fmt.Println("unknown user")
			return
		}
		// ユーザー名とパスワードの認証処理

		err = bcrypt.CompareHashAndPassword(ac[0].Password, []byte(password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
			return
		}

		if username == ac[0].Username {
			fmt.Printf("User %v logged in successfully\n", username)
			// 認証に成功した場合、セッションを開始し、ログイン済みの状態にする。
			session := sessions.Default(c)
			session.Set("username", username)
			session.Save()

			// ログインに成功した旨のメッセージを表示する。
			c.Redirect(http.StatusFound, "/view")
		} else {
			fmt.Printf("Login failure for user %v\n", username)
			// 認証に失敗した場合、ログインページにリダイレクトする。
			c.Redirect(http.StatusFound, "/login")
		}
	})

	r.GET("/logout", func(c *gin.Context) {
		// セッションを破棄し、ログアウトする。
		session := sessions.Default(c)
		username := session.Get("username")
		session.Clear()
		session.Save()
		fmt.Printf("%s logged out\n", username)
		c.Redirect(http.StatusFound, "/login")
	})

	r.Use(func(c *gin.Context) {
		// セッションからユーザー名を取得する。
		session := sessions.Default(c)
		username := session.Get("username")

		// ユーザー名が存在しない場合、ログインページにリダイレクトする。
		if username == nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}

		// ユーザー名が存在する場合、次のハンドラーに進む。
		c.Next()
	})

}
