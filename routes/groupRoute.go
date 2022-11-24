package routes

import (
	"net/http"
	"strconv"
	"strings"

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

	group := models.Group{Owner: user.ID, OwnerEmail: user.Email, GroupID: body.Id, Password: string(hash), Name: body.Name, OthersCanAdd: body.OthersCanAdd, OwnerName: user.Name, AboutGroup: body.AboutGroup, Particapants: strconv.FormatUint(uint64(user.ID), 10)}
	result := initializers.DB.Create(&group)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	groups := user.Groups + ":" + strconv.FormatUint(uint64(group.ID), 10)
	initializers.DB.Model(&models.User{}).Where("id = ?", user.ID).Update("groups", groups)

	c.JSON(http.StatusOK, gin.H{
		"groupName":    group.Name,
		"groupOwner":   group.OwnerName,
		"groupId":      group.GroupID,
		"aboutGroup":   group.AboutGroup,
		"particapants": group.Particapants,
	})
	return
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func GetGroupInfo(c *gin.Context) {
	var body struct {
		GroupId string `json:"groupId"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	user2, _ := c.Get("user")
	user := user2.(models.User)
	if body.GroupId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	var group models.Group
	initializers.DB.First(&group, "group_id = ?", body.GroupId)

	if group.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid GroupId",
		})
		return
	}

	users := strings.Split(group.Particapants, ":")
	if contains(users, strconv.FormatUint(uint64(user.ID), 10)) {
		var groupUsers []string
		for _, s := range users {
			var user models.User
			initializers.DB.First(&user, "id = ?", s)
			if user.ID != 0 {
				groupUsers = append(groupUsers, user.Name)
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"name":         group.Name,
			"owner":        group.OwnerName,
			"created":      group.CreatedAt,
			"group_id":     group.GroupID,
			"about_group":  group.AboutGroup,
			"particapants": groupUsers,
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"error": "User not in group.",
	})
	return
}
