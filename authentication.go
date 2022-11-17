package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
