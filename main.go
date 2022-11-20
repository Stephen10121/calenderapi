package main

import (
	"github.com/gin-gonic/gin"
	"github.com/stephen10121/calenderapi/initializers"
	"github.com/stephen10121/calenderapi/middleware"
	"github.com/stephen10121/calenderapi/routes"
)

//println(context.Query("test"))                       // Getting parameters test
//println(context.Request.Header.Get("Authorization")) // Getting header test

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectDatabase()
	initializers.SyncDatabase()
}

func main() {
	println("gello12")
	router := gin.Default()

	router.POST("/login", routes.Login)
	router.POST("/signup", routes.Signup)
	router.GET("/validate", middleware.RequireAuth, routes.Validate)

	router.Run()
}
