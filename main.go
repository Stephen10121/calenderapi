package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/stephen10121/calenderapi/initializers"
	"github.com/stephen10121/calenderapi/middleware"
	"github.com/stephen10121/calenderapi/routes"
)

func init() {
	// initializers.LoadEnvVariables()
	initializers.ConnectDatabase()
	initializers.SyncDatabase()
}

func main() {
	// gin.SetMode(gin.ReleaseMode) // Uncomment this in release
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	// Authentication

	// V1 Auth. Using Login/Signup.
	// router.POST("/login", routes.Login)
	// router.POST("/signup", routes.Signup)

	// V2 Auth. Using Google Authentication.
	router.POST("/google", routes.GoogleLogin)
	router.GET("/validate", middleware.RequireAuth, routes.Validate)
	router.POST("/addNotification", middleware.RequireAuth, routes.NotificationTokenAdd)

	// Group Part
	router.POST("/createGroup", middleware.RequireAuth, routes.CreateGroup)
	router.POST("/joinGroup", middleware.RequireAuth, routes.JoinGroup)
	router.POST("/leaveGroup", middleware.RequireAuth, routes.LeaveGroup)
	router.POST("/groupIdTaken", middleware.RequireAuth, routes.GroupIdTaken)
	router.POST("/groupInfo", middleware.RequireAuth, routes.GetGroupInfo)
	router.GET("/myGroups", middleware.RequireAuth, routes.GetMyGroups)
	router.POST("/acceptUser", middleware.RequireAuth, routes.AcceptParticapant)
	router.POST("/rejectUser", middleware.RequireAuth, routes.RejectParticapant)
	router.POST("/cancelRequest", middleware.RequireAuth, routes.CancelRequest)
	router.POST("/deleteGroup", middleware.RequireAuth, routes.RemoveGroup)
	router.POST("/kickUser", middleware.RequireAuth, routes.KickParticapant)

	// Job Part
	router.POST("/addJob", middleware.RequireAuth, routes.AddJob)
	router.POST("/getJobs", middleware.RequireAuth, routes.GetJobs)
	router.POST("/getAllJobsByMonthYear", middleware.RequireAuth, routes.GetJobsByMonthYear)
	router.POST("/allJobsByMonthsYear", middleware.RequireAuth, routes.GetJobsByMonthsYear)
	router.POST("/jobInfo", middleware.RequireAuth, routes.JobInfo)
	router.POST("/acceptJob", middleware.RequireAuth, routes.AcceptJob)
	fmt.Println("Running Server on ", 9090)
	router.Run()
}
