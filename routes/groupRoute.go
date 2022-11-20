package routes

import "github.com/gin-gonic/gin"

func CreateGroup(c *gin.Context) {
	var body struct {
		Name         string
		Id           string
		Password     string
		OthersCanAdd bool
	}
	println(body)
}
