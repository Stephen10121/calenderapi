package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stephen10121/calenderapi/initializers"
	"github.com/stephen10121/calenderapi/models"
)

type User2 struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	RName    string `json:"RName"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

type LoginReturn struct {
	Message   string `json:"message"`
	UserRName string `json:"usersRName"`
	AuthToken string `json:"authToken"`
}

var Users = []User2{
	{ID: "1", Username: "root", Password: "root", RName: "Stephen"},
}

func Login(context *gin.Context) {
	if context.Query("username") == "" || context.Query("password") == "" {
		context.IndentedJSON(http.StatusUnprocessableEntity, ErrorMessage{Message: "Missing parameters."})
		return
	}
	context.IndentedJSON(http.StatusOK, LoginReturn{Message: "Welcome", UserRName: Users[0].RName, AuthToken: "n21jkebdwkncewiu3r2i,32huidh"})
}

func Signup(context *gin.Context) {
	user := models.User{Username: "root", Password: "root"}

	result := initializers.DB.Create(&user)

	if result.Error != nil {
		context.Status(http.StatusBadRequest)
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
