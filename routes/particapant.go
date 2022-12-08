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
)

func RejectParticapant(c *gin.Context) {
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

	if userPart.ID == group.Owner {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "You cannot reject the owner.",
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

	realtime.UserGotRejected(group.ID, userPart.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Success.",
	})
	return
}

func KickParticapant(c *gin.Context) {
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

	if userPart.ID == group.Owner {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "You cannot kick out the owner.",
		})
		return
	}

	groupUsers := strings.Split(group.Particapants, ":")
	if functions.Contains(groupUsers, body.Particapant) != true {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not part of group anymore.",
		})
		return
	}

	usersGroups := strings.Split(userPart.Groups, ":")
	var groups string

	if len(usersGroups) != 0 {
		for i := 0; i < len(usersGroups); i++ {
			if usersGroups[i] != strconv.FormatUint(uint64(group.ID), 10) {
				groups = groups + ":" + usersGroups[i]
			}
		}
	}
	initializers.DB.Model(&models.User{}).Where("id = ?", userPart.ID).Update("groups", groups)

	groupParticapants := strings.Split(group.Particapants, ":")
	var particapants string

	if len(groupParticapants) != 0 {
		for i := 0; i < len(groupParticapants); i++ {
			if groupParticapants[i] != strconv.FormatUint(uint64(userPart.ID), 10) {
				particapants = particapants + ":" + groupParticapants[i]
			}
		}
	}
	initializers.DB.Model(&models.Group{}).Where("id = ?", group.ID).Update("particapants", particapants)

	realtime.UserKickedOut(group.ID, userPart.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Success.",
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

	realtime.UserGotAccepted(group.ID, userPart.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Success.",
	})
	return
}
