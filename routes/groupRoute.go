package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/stephen10121/calenderapi/functions"
	"github.com/stephen10121/calenderapi/initializers"
	"github.com/stephen10121/calenderapi/models"
	"github.com/stephen10121/calenderapi/realtime"
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
	if functions.Contains(users, strconv.FormatUint(uint64(user.ID), 10)) {
		var groupUsers []string
		for _, s := range users {
			var user models.User
			initializers.DB.First(&user, "id = ?", s)
			if user.ID != 0 {
				groupUsers = append(groupUsers, user.Name)
			}
		}

		if user.ID == group.Owner {
			usersPending := strings.Split(group.PendingParticapants, ":")
			var groupUsersPending []string
			for _, s := range usersPending {
				var user models.User
				initializers.DB.First(&user, "id = ?", s)
				if user.ID != 0 {
					groupUsersPending = append(groupUsersPending, user.Name)
				}
			}

			c.JSON(http.StatusOK, gin.H{
				"name":                 group.Name,
				"owner":                group.OwnerName,
				"created":              group.CreatedAt,
				"group_id":             group.GroupID,
				"about_group":          group.AboutGroup,
				"particapants":         groupUsers,
				"pending_particapants": groupUsersPending,
			})
			return
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

func JoinGroup(c *gin.Context) {
	var body struct {
		Id       string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	user2, _ := c.Get("user")
	user := user2.(models.User)

	if body.Id == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	var group models.Group
	initializers.DB.First(&group, "group_id = ?", body.Id)

	if group.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Id or password",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(group.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Id or password",
		})
		return
	}

	usersCurrentPendingGroups := strings.Split(user.PendingGroups, ":")
	if functions.Contains(usersCurrentPendingGroups, strconv.FormatUint(uint64(group.ID), 10)) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Already joined or attempted to join",
		})
		return
	}

	usersCurrentGroups := strings.Split(user.Groups, ":")
	if functions.Contains(usersCurrentGroups, strconv.FormatUint(uint64(group.ID), 10)) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Already joined or attempted to join",
		})
		return
	}

	groups := user.PendingGroups + ":" + strconv.FormatUint(uint64(group.ID), 10)
	initializers.DB.Model(&models.User{}).Where("id = ?", user.ID).Update("pending_groups", groups)

	particapants := group.PendingParticapants + ":" + strconv.FormatUint(uint64(user.ID), 10)
	initializers.DB.Model(&models.Group{}).Where("id = ?", group.ID).Update("pending_particapants", particapants)

	messageToOwner := user.Name + " wants to join your group."
	realtime.NotifyGroupOwner(group.ID, "Pending New User", messageToOwner)

	c.JSON(http.StatusOK, gin.H{
		"error":   "none",
		"message": "Success. Now wait for the group owner to accept the join request.",
	})
	return
}

type groupPendingData struct {
	GroupId   string `json:"groupId"`
	GroupName string `json:"groupName"`
}

type groupData struct {
	GroupId      string `json:"groupId"`
	GroupName    string `json:"groupName"`
	GroupOwner   string `json:"groupOwner"`
	OthersCanAdd bool   `json:"othersCanAdd"`
}

func GetMyGroups(c *gin.Context) {
	user2, _ := c.Get("user")
	user := user2.(models.User)

	usersCurrentPendingGroups := strings.Split(user.PendingGroups, ":")
	var userCurrentPendingGroupsJson []groupPendingData
	if len(usersCurrentPendingGroups) != 0 {
		for i := 0; i < len(usersCurrentPendingGroups); i++ {
			var group models.Group
			u64, err := strconv.ParseUint(usersCurrentPendingGroups[i], 10, 16)
			if err != nil {
				fmt.Println(err)
				continue
			}
			initializers.DB.First(&group, u64)

			if group.ID == 0 {
				continue
			}
			userCurrentPendingGroupsJson = append(userCurrentPendingGroupsJson, groupPendingData{GroupId: group.GroupID, GroupName: group.Name})
		}
	}
	usersCurrentGroups := strings.Split(user.Groups, ":")
	var usersCurrentGroupsJson []groupData

	if len(usersCurrentGroups) != 0 {
		for i := 0; i < len(usersCurrentGroups); i++ {
			var group models.Group
			u64, err := strconv.ParseUint(usersCurrentGroups[i], 10, 16)
			if err != nil {
				fmt.Println(err)
				continue
			}
			initializers.DB.First(&group, u64)

			if group.ID == 0 {
				continue
			}
			usersCurrentGroupsJson = append(usersCurrentGroupsJson, groupData{GroupId: group.GroupID, GroupName: group.Name, OthersCanAdd: group.OthersCanAdd, GroupOwner: group.OwnerName})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"groups":        usersCurrentGroupsJson,
		"pendingGroups": userCurrentPendingGroupsJson,
	})
	return
}

func AcceptParticapant(c *gin.Context) {
	var body struct {
		Id          string `json:"id"`
		Particapant string `json:"particapant"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	user2, _ := c.Get("user")
	user := user2.(models.User)

	if body.Id == "" || body.Particapant == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	var group models.Group
	initializers.DB.First(&group, "group_id = ?", body.Id)

	if group.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Group",
		})
		return
	}

	if user.ID != group.Owner {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Your not the owner.",
		})
		return
	}

	var userPart models.User
	u64, err := strconv.ParseUint(body.Particapant, 10, 16)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid User Id",
		})
		return
	}
	initializers.DB.First(&userPart, u64)

	if userPart.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid User",
		})
		return
	}

	groupPendingUsers := strings.Split(group.PendingParticapants, ":")
	if functions.Contains(groupPendingUsers, body.Particapant) != true {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not pending anymore.",
		})
		return
	}

	usersPendingGroups := strings.Split(userPart.PendingGroups, ":")
	var pendingGroups string

	if len(usersPendingGroups) != 0 {
		for i := 0; i < len(usersPendingGroups); i++ {
			if usersPendingGroups[i] != strconv.FormatUint(uint64(group.ID), 10) {
				pendingGroups = pendingGroups + ":" + usersPendingGroups[i]
			}
		}
	}
	initializers.DB.Model(&models.User{}).Where("id = ?", userPart.ID).Update("pending_groups", pendingGroups)
	groups := userPart.Groups + ":" + strconv.FormatUint(uint64(group.ID), 10)
	initializers.DB.Model(&models.User{}).Where("id = ?", userPart.ID).Update("groups", groups)

	groupPendingParticapants := strings.Split(group.PendingParticapants, ":")
	var pendingParticapants string

	if len(groupPendingParticapants) != 0 {
		for i := 0; i < len(groupPendingParticapants); i++ {
			if groupPendingParticapants[i] != strconv.FormatUint(uint64(userPart.ID), 10) {
				pendingParticapants = pendingParticapants + ":" + groupPendingParticapants[i]
			}
		}
	}
	initializers.DB.Model(&models.Group{}).Where("id = ?", group.ID).Update("pending_particapants", pendingParticapants)
	users := group.Particapants + ":" + strconv.FormatUint(uint64(userPart.ID), 10)
	initializers.DB.Model(&models.Group{}).Where("id = ?", group.ID).Update("particapants", users)

	c.JSON(http.StatusOK, gin.H{
		"error":   "none",
		"message": "Success.",
	})
	return
}
