package main

import (
	"github.com/gin-gonic/gin"
)

//println(context.Query("test"))                       // Getting parameters test
//println(context.Request.Header.Get("Authorization")) // Getting header test

func init() {
	LoadEnvVariables()
	ConnectDatabase()
}

func main() {
	println("gello12")
	router := gin.Default()

	router.GET("/login", Login)

	router.Run()
}
