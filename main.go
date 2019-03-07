package main

import (
	"log"
	"net/http"
)

var sp = "xuehua5201314"
var secret = "8fbd540d0b23197df1d5095f0d6ee46d"
var appId = "wxa2c324b63b2a9e5e"
var gameUrl = "https://wxwyjh.chiji-h5.com"

func main() {
	http.HandleFunc("/upload", UploadHandler)
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
