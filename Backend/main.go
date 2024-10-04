package main

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"regexp"
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
	writeResult(c, http.StatusOK, []byte(`{"name": "Scott"}`))
}

func handleLogin(c *gin.Context) {
	var user User
	user.username = c.PostForm("username")
	password := c.PostForm("password")
	err := user.getUserByUserName()
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

	session, _ := store.Get(c.Request, "session")
	session.Values["id"] = user
	session.Save(c.Request, c.Writer)
	writeResult(c, http.StatusOK, []byte(`{"success": "User successfully logged in"}`))
}

func handleRegister(c *gin.Context) {
	var user User

	if err := c.BindJSON(&user); err != nil {
		writeResult(c, http.StatusBadRequest, []byte(`{"error": "Bad request"}`))
		return
	}

	if !validateEmail(user.email) || !validateUserName(user.username) {
		writeResult(c, http.StatusBadRequest, []byte(`{"error": "Invalid Email or username"}`))
		return
	}

	query := "INSERT INTO Users (Email, PassCode, SettingsBox, SettingsGuess) VALUES (?, ?, ?, ?)"
	res, err := db.Exec(query, user.email, user.pw, user.settings_box, user.settings_guess)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted into DB: ", res)
}

func validateEmail(s string) bool {
	_, err := mail.ParseAddress(s)
	return err == nil
}

func validateUserName(s string) bool {
	re := regexp.MustCompile("^(?=.*[A-Za-z])(?=.*[0-9])(?=.*)[A-Za-z0-9]{8,20}$*/")
	return re.MatchString(s)
}

func (u *User) getUserByUserName() error {
	query := "SELECT * FROM Users WHERE Username = ?"

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
	router.GET("/register", handleRegister)
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
	db, err = sql.Open("mysql", "root:secret-jungle@tcp(localhost:3306)/GuessingGameDb")
	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	initDb()
	r := setupRouter()
	r.Run("localhost:3000")
	defer db.Close()
}
