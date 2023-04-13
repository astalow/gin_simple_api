package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type User struct {
	Username string
	Password string
}

func session(r *gin.Engine) {
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	r.POST("/login", func(c *gin.Context) {
		// RequestLogger(c)
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			// JSONパースエラーが発生した場合
			fmt.Println("nande")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}
		fmt.Println(user.Username)
		fmt.Println(user.Password)

		// フォームから送信されたユーザー名とパスワードを取得する。
		username := user.Username
		password := user.Password

		// ユーザー名とパスワードの認証処理を実装する。
		// ここでは、簡単のために固定のユーザー名とパスワードを設定しています。
		if username == "user" && password == "password" {
			fmt.Println("login success")
			// 認証に成功した場合、セッションを開始し、ログイン済みの状態にする。
			session := sessions.Default(c)
			session.Set("username", username)
			session.Save()

			// ログインに成功した旨のメッセージを表示する。
			c.Redirect(http.StatusFound, "/view")
		} else {
			fmt.Println("login failure")
			// 認証に失敗した場合、ログインページにリダイレクトする。
			c.Redirect(http.StatusFound, "/login")
		}
	})

	r.GET("/logout", func(c *gin.Context) {
		// セッションを破棄し、ログアウトする。
		session := sessions.Default(c)
		session.Clear()
		session.Save()

		// ログアウトした旨のメッセージを表示する。
		c.HTML(http.StatusOK, "logout.html", nil)
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

	r.GET("/secret", func(c *gin.Context) {
		// セッションからユーザー名を取得する。
		session := sessions.Default(c)
		username := session.Get("username")

		// ユーザー名を表示する。
		c.HTML(http.StatusOK, "secret.html", gin.H{"username": username})
	})

}
