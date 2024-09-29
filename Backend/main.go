package main

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	id             uint32
	email          string
	username       string
	pw             string
	settings_guess int
	settings_box   int
}

var db *sql.DB

var store = sessions.NewCookieStore([]byte("secret-jungle"))

func init() {
	store.Options.HttpOnly = true
	// store.Options.Secure = true // https
	gob.Register(&User{})
}

func writeResult(c *gin.Context, status int, data []byte) {
	c.Data(status, "application/json", data)
}

func authHandler(c *gin.Context) {
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

func serveIndex(c *gin.Context) {
	fmt.Println("HELLO FROM GO")
	writeResult(c, http.StatusOK, []byte(`{"name": "Scott"}`))
}

func handleLogin(c *gin.Context) {
	fmt.Println("IN HANDLER")
	var user User
	user.username = c.PostForm("username")
	password := c.PostForm("password")
	err := user.getUserByUserName()
	fmt.Println("AFTER GET USERNAME")
	if err != nil {
		fmt.Println("error retrieving user from DB ", err)
		writeResult(c, http.StatusUnauthorized, []byte(`{"error": "Incorrect username or password"}`))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.pw), []byte(password))

	if err != nil {
		fmt.Println("error with password, ", err)
		writeResult(c, http.StatusUnauthorized, []byte(`{"error": "Incorrect username or password"}`))
		return
	}
	fmt.Println("PASSED BCRYPT CHECK")

	session, _ := store.Get(c.Request, "session")
	session.Values["id"] = user
	session.Save(c.Request, c.Writer)
	writeResult(c, http.StatusOK, []byte(`{"success": "User successfully logged in"}`))
}

func (u *User) getUserByUserName() error {
	query := "SELECT * FROM Users WHERE Username = ?"
	fmt.Println("BEFORE ROW CHECK ", db.Stats().OpenConnections)
	if err2 := db.Ping(); err2 != nil {
		fmt.Println("PING ERR ", err2)
	}
	fmt.Println("AFTER ROW CHECK ")
	if err := db.QueryRow(query, u.username).Scan(&u.id, &u.email, &u.pw, &u.settings_guess, &u.settings_box); err != nil {
		fmt.Println("getUserByUserName() returned error: ", err)
		return err
	}
	return nil
}

func setupRouter() *gin.Engine {
	router := gin.Default() // init router with default mw (e.g. logging)

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "PUT", "PATCH", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))

	router.GET("/", serveIndex)
	router.GET("/usersLoggedIn")
	router.GET("/register")
	router.POST("/login", handleLogin)
	router.GET("/logout")

	user := router.Group("/user", authHandler) // User specific routes
	{
		user.GET("/game")
		user.GET("/stats")
		user.GET("/history")
		user.GET("/updateUser")
	}

	return router
}

func initDb() {
	var err error
	db, err = sql.Open("mysql", "root:secret-jungle@tcp(localhost:3000)/GuessingGameDb")
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	initDb()
	r := setupRouter()
	r.Run("localhost:3000")
	defer db.Close()
}
