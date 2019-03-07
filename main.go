package main

import (
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"log"
)

var secret = "8fbd540d0b23197df1d5095f0d6ee46d"
var appId = "wxa2c324b63b2a9e5e"
var gameUrl = "https://wxwyjh.chiji-h5.com"

func main() {
	router := gin.Default()
	router.GET("/upload", UploadHandler)
	if err := endless.ListenAndServe(":8888", router); err != nil {
		log.Fatal("ListenAndServe-error: ", err)
	}
}
