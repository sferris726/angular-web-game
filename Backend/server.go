package main

import (
	//"net/http"
	"github.com/gin-gonic/gin"
	"fmt"
)

type ConInfo struct {
	host     string
	user     string
	pw       string
	database string
}

var con_info = ConInfo{host: "localhost", user: "", pw: "", database: "GuessingGameDb"}

func serveIndex(context *gin.Context) {
	fmt.Println("HELLO WORLD")
}

func deployServer() {
	router := gin.Default()
	router.GET("/", serveIndex)
	router.GET("/usersLoggedIn")
	router.Run("localhost:3000")
}
