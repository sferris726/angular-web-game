package main

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

type User struct {
	id             uint32
	email          string
	pw             string
	settings_guess int
	settings_nox   int
}

var db *sql.DB

var store = sessions.NewCookieStore([]byte("secret-jungle"))

func init() {
	store.Options.HttpOnly = true
	store.Options.Secure = true // https
	gob.Register(&User{})
}

func writeResult(c *gin.Context, status int, data []byte) {
	c.Data(status, "application/json", data)
}

func auth(c *gin.Context) {
	session, _ := store.Get(c.Request, "session")
	fmt.Println("session", session)
	_, ok := session.Values["user"]

	if !ok {
		writeResult(c, http.StatusForbidden, []byte(`{"error": "bad session"}`))
		c.Abort()
		return
	}

	c.Next()
}

func setupRouter() *gin.Engine {
	router := gin.Default() // init router with default mw (e.g. logging)
	router.GET("/usersLoggedIn")
	router.GET("/register")
	router.GET("/login")
	router.GET("/logout")
	router.GET("/game")
	router.GET("/stats")
	router.GET("/history")
	router.GET("/updateUser")
	// user := router.Group("/api/user")
	// {
	// 	user.POST("/", controllers.CreateUser)
	// 	user.GET("/:userId", controllers.GetUserById)
	// 	user.GET("/", controllers.GetAllUsers)
	// }

	return router
}

func initDb() {
	var err error
	db, err = sql.Open("mysql", "root:secret-jungle@tcp(localhost:3000)/gin_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
}

func main() {
	r := setupRouter()
	initDb()

	authRouter := r.Group("/user")

	r.Run("localhost:3000")
}
