package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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

	var groupPendingParticapants []uint
	json.Unmarshal([]byte(group.PendingParticapants), &groupPendingParticapants)
	if functions.UintContains(groupPendingParticapants, userPart.ID) != true {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not pending anymore.",
		})
		return
	}

	var pendingGroups []uint
	var userPartPendingGroups []uint
	json.Unmarshal([]byte(userPart.PendingGroups), &userPartPendingGroups)
	for _, s := range userPartPendingGroups {
		if s != group.ID {
			pendingGroups = append(pendingGroups, s)
		}
	}
	userPartPendingGroupsJson, _ := json.Marshal(pendingGroups)
	initializers.DB.Model(&models.User{}).Where("id = ?", userPart.ID).Update("pending_groups", userPartPendingGroupsJson)

	var pendingParticapants []uint
	var groupsPendingParticapants []uint
	json.Unmarshal([]byte(group.PendingParticapants), &groupsPendingParticapants)
	for _, s := range groupsPendingParticapants {
		if s != userPart.ID {
			pendingParticapants = append(pendingParticapants, s)
		}
	}
	groupsPendingParticapantsJson, _ := json.Marshal(pendingParticapants)
	initializers.DB.Model(&models.Group{}).Where("id = ?", group.ID).Update("pending_particapants", groupsPendingParticapantsJson)

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

	var groupParticapants []uint
	json.Unmarshal([]byte(group.Particapants), &groupParticapants)
	if functions.UintContains(groupParticapants, userPart.ID) != true {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not part of group anymore.",
		})
		return
	}

	var groups []uint
	var userPartGroups []uint
	json.Unmarshal([]byte(userPart.Groups), &userPartGroups)
	for _, s := range userPartGroups {
		if s != group.ID {
			groups = append(groups, s)
		}
	}
	userPartGroupsJson, _ := json.Marshal(groups)
	initializers.DB.Model(&models.User{}).Where("id = ?", userPart.ID).Update("groups", userPartGroupsJson)

	var particapants []uint
	var groupsParticapants []uint
	json.Unmarshal([]byte(group.Particapants), &groupsParticapants)
	for _, s := range groupsParticapants {
		if s != userPart.ID {
			particapants = append(particapants, s)
		}
	}
	groupsParticapantsJson, _ := json.Marshal(particapants)
	initializers.DB.Model(&models.Group{}).Where("id = ?", group.ID).Update("particapants", groupsParticapantsJson)

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

	var groupPendingParticapants []uint
	json.Unmarshal([]byte(group.PendingParticapants), &groupPendingParticapants)
	if functions.UintContains(groupPendingParticapants, userPart.ID) != true {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not pending anymore.",
		})
		return
	}

	var pendingGroups []uint
	var userPartPendingGroups []uint
	json.Unmarshal([]byte(userPart.PendingGroups), &userPartPendingGroups)
	for _, s := range userPartPendingGroups {
		if s != group.ID {
			pendingGroups = append(pendingGroups, s)
		}
	}
	userPartPendingGroupsJson, _ := json.Marshal(pendingGroups)
	initializers.DB.Model(&models.User{}).Where("id = ?", userPart.ID).Update("pending_groups", userPartPendingGroupsJson)

	var groups []uint
	json.Unmarshal([]byte(userPart.Groups), &groups)
	groups = append(groups, group.ID)
	groupsJson, _ := json.Marshal(groups)
	initializers.DB.Model(&models.User{}).Where("id = ?", userPart.ID).Update("groups", groupsJson)

	var pendingParticapants []uint
	var groupsPendingParticapants []uint
	json.Unmarshal([]byte(group.PendingParticapants), &groupsPendingParticapants)
	for _, s := range groupsPendingParticapants {
		if s != userPart.ID {
			pendingParticapants = append(pendingParticapants, s)
		}
	}
	groupsPendingParticapantsJson, _ := json.Marshal(pendingParticapants)
	initializers.DB.Model(&models.Group{}).Where("id = ?", group.ID).Update("pending_particapants", groupsPendingParticapantsJson)

	var users []uint
	json.Unmarshal([]byte(group.Particapants), &users)
	users = append(users, userPart.ID)
	usersJson, _ := json.Marshal(users)
	initializers.DB.Model(&models.Group{}).Where("id = ?", group.ID).Update("particapants", usersJson)

	realtime.UserGotAccepted(group.ID, userPart.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Success.",
	})
	return
}
