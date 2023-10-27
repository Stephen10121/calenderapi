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

func LeaveGroup(c *gin.Context) {
	var body struct {
		Id       string `json:"id"`
		Transfer string `json:"transfer"`
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
		if body.Transfer == "0" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "New owner was not selected.",
			})
			return
		}

		var userInPart models.User
		initializers.DB.First(&userInPart, body.Transfer)

		if userInPart.ID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "New owner doesnt exist.",
			})
			return
		}

		groupParticapants, err := helpers.UnmarshalGroupParticapants(group.Particapants)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Something went wrong.",
			})
			fmt.Println(err)
			return
		}

		if functions.UintContains(groupParticapants, userInPart.ID) != true {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "User not part of the group.",
			})
			return
		}

		var groups []uint
		var userGroups []uint
		json.Unmarshal([]byte(user.Groups), &userGroups)
		for _, s := range userGroups {
			if s != group.ID {
				groups = append(groups, s)
			}
		}
		groupsJson, _ := json.Marshal(groups)
		initializers.DB.Model(&models.User{}).Where("id = ?", user.ID).Update("groups", groupsJson)

		var particapants []uint
		var groupsParticapants []uint
		json.Unmarshal([]byte(group.Particapants), &groupsParticapants)
		for _, s := range groupsParticapants {
			if s != user.ID {
				particapants = append(particapants, s)
			}
		}
		particapantsJson, _ := json.Marshal(particapants)
		initializers.DB.Model(&models.Group{}).Where("id = ?", group.ID).Update("particapants", particapantsJson)

		initializers.DB.Model(&models.Group{}).Where("id = ?", group.ID).Update("owner", userInPart.ID)
		initializers.DB.Model(&models.Group{}).Where("id = ?", group.ID).Update("owner_email", userInPart.Email)
		initializers.DB.Model(&models.Group{}).Where("id = ?", group.ID).Update("owner_name", userInPart.FullName)

		realtime.UserLeftTransfered(group.GroupID, string(particapantsJson), userInPart.FullName)

		c.JSON(http.StatusOK, gin.H{
			"message": "Success.",
		})
		return
	}

	var groupsParticapants []uint
	json.Unmarshal([]byte(group.Particapants), &groupsParticapants)
	if functions.UintContains(groupsParticapants, user.ID) != true {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Your not part of the group.",
		})
		return
	}

	var groups []uint
	var userGroups []uint
	json.Unmarshal([]byte(user.Groups), &userGroups)
	for _, s := range userGroups {
		if s != group.ID {
			groups = append(groups, s)
		}
	}
	groupsJson, _ := json.Marshal(groups)
	initializers.DB.Model(&models.User{}).Where("id = ?", user.ID).Update("groups", groupsJson)

	var particapants []uint

	groupParticapants, err := helpers.UnmarshalGroupParticapants(group.Particapants)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong.",
		})
		fmt.Println(err)
		return
	}

	for _, s := range groupParticapants {
		if s != user.ID {
			particapants = append(particapants, s)
		}
	}
	groupParticapantsJson, _ := json.Marshal(particapants)
	initializers.DB.Model(&models.Group{}).Where("id = ?", group.ID).Update("particapants", groupParticapantsJson)

	realtime.UserLeft(group.ID, user.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Success.",
	})
	return
}
