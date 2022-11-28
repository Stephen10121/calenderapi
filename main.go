package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stephen10121/calenderapi/initializers"
	"github.com/stephen10121/calenderapi/middleware"
	"github.com/stephen10121/calenderapi/routes"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectDatabase()
	initializers.SyncDatabase()
}

func main() {
	// gin.SetMode(gin.ReleaseMode) // Uncomment this in release
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())
	// Authentication
	router.POST("/login", routes.Login)
	router.POST("/signup", routes.Signup)
	router.GET("/validate", middleware.RequireAuth, routes.Validate)
	// Group Part
	router.POST("/createGroup", middleware.RequireAuth, routes.CreateGroup)
	router.POST("/joinGroup", middleware.RequireAuth, routes.JoinGroup)
	router.POST("/groupIdTaken", middleware.RequireAuth, routes.GroupIdTaken)
	router.POST("/groupInfo", middleware.RequireAuth, routes.GetGroupInfo)
	router.GET("/myGroups", middleware.RequireAuth, routes.GetMyGroups)
	router.POST("/acceptUser", middleware.RequireAuth, routes.AcceptParticapant)
	// Job Part
	router.POST("/addJob", middleware.RequireAuth, routes.AddJob)
	router.POST("/getJobs", middleware.RequireAuth, routes.GetJobs)
	fmt.Println("Running Server on ", os.Getenv("PORT"))
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
