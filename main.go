package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

var secret = "8fbd540d0b23197df1d5095f0d6ee46d"
var appId = "wxa2c324b63b2a9e5e"
var gameUrl = "https://wxwyjh.chiji-h5.com"

func main() {
	engine := gin.Default()
	engine.GET("/upload", UploadHandler)
	if err := engine.Run(":8888"); err != nil {
		log.Fatal("ListenAndServe-error: ", err)
	}
}
