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
	Email         string `json:"email"`
	Username      string `json:"username"`
	Pw            string `json:"password"`
	SettingsGuess int    `json:"settings_guess"`
	SettingsBox   int    `json:"settings_box"`
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
	user.Username = c.PostForm("username")
	password := c.PostForm("password")
	err := user.getUserByUserName()
	if err != nil {
		fmt.Println("error retrieving user from DB ", err)
		writeResult(c, http.StatusUnauthorized, []byte(`{"error": "Incorrect username or password"}`))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Pw), []byte(password))

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
	fmt.Println("MADE IT IN")
	if err := c.BindJSON(&user); err != nil {
		fmt.Println("BAD REQ")
		writeResult(c, http.StatusBadRequest, []byte(`{"error": "Bad request"}`))
		return
	}

	if !validateEmail(user.Email) /* || !validateUserName(user.username)*/ {
		fmt.Println("INVALID EMAIL: ", user.Email)
		writeResult(c, http.StatusBadRequest, []byte(`{"error": "Invalid Email or username"}`))
		return
	}

	query := "INSERT INTO Users (Email, Username, PassCode, SettingsBox, SettingsGuess) VALUES (?, ?, ?, ?, ?)"
	_, err := db.Exec(query, user.Email, user.Username, user.Pw, user.SettingsGuess, user.SettingsBox)
	if err != nil {
		log.Fatal(err)
	}

	writeResult(c, http.StatusOK, []byte(`{"success": "User registered successfully!"}`))
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

	if err := db.QueryRow(query, u.Username).Scan(&u.Email, &u.Pw, &u.SettingsBox, &u.SettingsGuess); err != nil {
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
	router.POST("/register", handleRegister)
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
