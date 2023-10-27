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
)

type PartiacapantSend struct {
	Name string `json:"name"`
	Id   uint   `json:"id"`
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

	groupParticapants, err := helpers.UnmarshalGroupParticapants(group.Particapants)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong.",
		})
		fmt.Println(err)
		return
	}

	if functions.UintContains(groupParticapants, user.ID) != true {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not in group.",
		})
		return
	}

	var groupUsers []PartiacapantSend
	for _, s := range groupParticapants {
		var user models.User
		initializers.DB.First(&user, "id = ?", s)
		if user.ID != 0 {
			groupUsers = append(groupUsers, PartiacapantSend{Name: user.FirstName + " " + user.LastName, Id: user.ID})
		}
	}

	if user.ID == group.Owner {
		var groupUsersPending []PartiacapantSend

		fmt.Println(group.PendingParticapants)
		var groupPendingParticapants []uint
		json.Unmarshal([]byte(group.PendingParticapants), &groupPendingParticapants)
		for _, s := range groupPendingParticapants {
			var user2 models.User
			fmt.Println(s)
			initializers.DB.First(&user2, "id = ?", s)
			if user2.ID != 0 {
				groupUsersPending = append(groupUsersPending, PartiacapantSend{Name: user2.FirstName + " " + user2.LastName, Id: user2.ID})
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"name":         group.Name,
			"owner":        group.OwnerName,
			"owner_email":  group.OwnerEmail,
			"created":      group.CreatedAt,
			"group_id":     group.GroupID,
			"about_group":  group.AboutGroup,
			"particapants": groupUsers,
			"yourowner": gin.H{
				"ownerId":              group.Owner,
				"pending_particapants": groupUsersPending,
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name":         group.Name,
		"owner_email":  group.OwnerEmail,
		"owner":        group.OwnerName,
		"created":      group.CreatedAt,
		"group_id":     group.GroupID,
		"about_group":  group.AboutGroup,
		"particapants": groupUsers,
		"yourowner":    false,
	})
	return
}
