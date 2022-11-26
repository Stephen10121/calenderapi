package main

import (
	"github.com/gin-gonic/gin"
	"github.com/stephen10121/calenderapi/initializers"
	"github.com/stephen10121/calenderapi/middleware"
	"github.com/stephen10121/calenderapi/routes"
)

//println(context.Request.Header.Get("Authorization")) // Getting header test

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectDatabase()
	initializers.SyncDatabase()
}

func main() {
	router := gin.Default()

	router.POST("/login", routes.Login)
	router.POST("/signup", routes.Signup)
	router.GET("/validate", middleware.RequireAuth, routes.Validate)
	router.POST("/createGroup", middleware.RequireAuth, routes.CreateGroup)
	router.POST("/joinGroup", middleware.RequireAuth, routes.JoinGroup)
	router.POST("/groupIdTaken", middleware.RequireAuth, routes.GroupIdTaken)
	router.POST("/addJob", middleware.RequireAuth, routes.AddJob)
	router.POST("/groupInfo", middleware.RequireAuth, routes.GetGroupInfo)

	// Get groups for each user. The user will contain an array with the group id their in.
	//

	router.Run()
}

// The react native fetch function.
//fetch('https://mywebsite.com/endpoint/', {
//	method: 'POST',
//	headers: {
//	  Accept: 'application/json',
//	  'Content-Type': 'application/json'
//	},
//	body: JSON.stringify({
//	  firstParam: 'yourValue',
//	  secondParam: 'yourOtherValue'
//	})
//  });
