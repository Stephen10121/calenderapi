package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stephen10121/calenderapi/initializers"
	"github.com/stephen10121/calenderapi/models"
	"golang.org/x/crypto/bcrypt"
)

func GroupIdTaken(c *gin.Context) {
	var body struct {
		Id string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	if len(body.Id) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	var group models.Group
	r := initializers.DB.Where(&models.Group{GroupID: body.Id}).First(&group)

	exists := r.RowsAffected > 0

	if exists {
		c.JSON(http.StatusOK, gin.H{
			"taken": true,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"taken": false,
	})
	return
}

func CreateGroup(c *gin.Context) {
	var body struct {
		Name         string
		Id           string
		Password     string
		OthersCanAdd bool
		AboutGroup   string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	user2, _ := c.Get("user")
	user := user2.(models.User)

	if body.Name == "" || body.Id == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	group := models.Group{Owner: user.ID, OwnerEmail: user.Email, GroupID: body.Id, Password: string(hash), Name: body.Name, OthersCanAdd: body.OthersCanAdd, OwnerName: user.Name, AboutGroup: body.AboutGroup}
	result := initializers.DB.Create(&group)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"groupName":  group.Name,
		"groupOwner": group.OwnerName,
		"groupId":    group.GroupID,
		"aboutGroup": group.AboutGroup,
	})
	return
}
