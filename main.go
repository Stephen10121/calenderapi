package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	grouproutes "github.com/stephen10121/calenderapi/groupRoutes"
	"github.com/stephen10121/calenderapi/initializers"
	"github.com/stephen10121/calenderapi/middleware"
	"github.com/stephen10121/calenderapi/routes"
	"github.com/stephen10121/calenderapi/variables"
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
	router.GET("/", routes.HomePage)
	router.POST("/google", routes.GoogleLogin)
	router.GET("/validate", middleware.RequireAuth, routes.Validate)
	router.POST("/addNotification", middleware.RequireAuth, routes.NotificationTokenAdd)

	// Group Part
	router.GET("/myGroups", middleware.RequireAuth, grouproutes.GetMyGroups)
	router.POST("/createGroup", middleware.RequireAuth, grouproutes.CreateGroup)
	router.POST("/joinGroup", middleware.RequireAuth, grouproutes.JoinGroup)
	router.POST("/leaveGroup", middleware.RequireAuth, grouproutes.LeaveGroup)
	router.POST("/groupIdTaken", middleware.RequireAuth, grouproutes.GroupIdTaken)
	router.POST("/groupInfo", middleware.RequireAuth, grouproutes.GetGroupInfo)
	router.POST("/acceptUser", middleware.RequireAuth, routes.AcceptParticapant)
	router.POST("/rejectUser", middleware.RequireAuth, routes.RejectParticapant)
	router.POST("/cancelRequest", middleware.RequireAuth, grouproutes.CancelRequest)
	router.POST("/deleteGroup", middleware.RequireAuth, grouproutes.RemoveGroup)
	router.POST("/kickUser", middleware.RequireAuth, routes.KickParticapant)

	// Job Part
	router.POST("/addJob", middleware.RequireAuth, routes.AddJob)
	router.POST("/getJobs", middleware.RequireAuth, routes.GetJobs)
	router.POST("/getAllJobsByMonthYear", middleware.RequireAuth, routes.GetJobsByMonthYear)
	router.POST("/allJobsByMonthsYear", middleware.RequireAuth, routes.GetJobsByMonthsYear)
	router.POST("/jobInfo", middleware.RequireAuth, routes.JobInfo)
	router.POST("/acceptJob", middleware.RequireAuth, routes.AcceptJob)

	fmt.Println("Running Server on ", variables.Port())
	router.Run(variables.Port())
}
