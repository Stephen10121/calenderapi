package grouproutes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stephen10121/calenderapi/functions"
	"github.com/stephen10121/calenderapi/helpers"
	"github.com/stephen10121/calenderapi/initializers"
	"github.com/stephen10121/calenderapi/models"
	"github.com/stephen10121/calenderapi/realtime"
)

func CancelRequest(c *gin.Context) {
	var body struct {
		Id string `json:"id"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	user2, _ := c.Get("user")
	user := user2.(models.User)

	if body.Id == "" {
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

	if user.ID == group.Owner {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Your the owner.",
		})
		return
	}

	var groupParticapants []uint
	json.Unmarshal([]byte(group.PendingParticapants), &groupParticapants)
	if functions.UintContains(groupParticapants, user.ID) != true {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Your not pending.",
		})
		return
	}

	var groups []uint
	userPendingGroups, err := helpers.UnmarshalPendingGroups(user.PendingGroups)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong",
		})
		fmt.Println(err)
		return
	}
	for _, s := range userPendingGroups {
		if s != group.ID {
			groups = append(groups, s)
		}
	}
	groupsJson, _ := json.Marshal(groups)
	initializers.DB.Model(&models.User{}).Where("id = ?", user.ID).Update("pending_groups", groupsJson)

	var particapants []uint
	for _, s := range groupParticapants {
		if s != user.ID {
			particapants = append(particapants, s)
		}
	}
	particapantsJson, _ := json.Marshal(particapants)
	initializers.DB.Model(&models.Group{}).Where("id = ?", group.ID).Update("pending_particapants", particapantsJson)

	realtime.UserLeftWhilePending(group.ID, user.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Success.",
	})
	return
}
