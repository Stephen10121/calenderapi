package grouproutes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stephen10121/calenderapi/helpers"
	"github.com/stephen10121/calenderapi/initializers"
	"github.com/stephen10121/calenderapi/models"
	"github.com/stephen10121/calenderapi/realtime"
)

func RemoveGroup(c *gin.Context) {
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

	if user.ID != group.Owner {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Your not the owner.",
		})
		return
	}

	var groupParticapants []uint
	json.Unmarshal([]byte(group.Particapants), &groupParticapants)
	for _, s := range groupParticapants {
		var userInPart models.User
		initializers.DB.First(&userInPart, "id = ?", s)

		if userInPart.ID == 0 {
			continue
		}

		var newGroups []uint
		var userInPartGroups []uint
		json.Unmarshal([]byte(userInPart.Groups), &userInPartGroups)
		for _, s2 := range userInPartGroups {
			if s2 != group.ID {
				newGroups = append(newGroups, s2)
			}
		}
		userInPartGroupsJson, _ := json.Marshal(newGroups)
		initializers.DB.Model(&models.User{}).Where("id = ?", s).Update("groups", userInPartGroupsJson)
	}

	var groupPendingParticapants []uint
	json.Unmarshal([]byte(group.PendingParticapants), &groupPendingParticapants)
	for _, s := range groupPendingParticapants {
		var userInPart models.User
		initializers.DB.First(&userInPart, "id = ?", s)

		if userInPart.ID == 0 {
			continue
		}

		var newGroups []uint
		userInPartPendingGroups, err := helpers.UnmarshalPendingGroups(user.PendingGroups)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Something went wrong",
			})
			fmt.Println(err)
			return
		}
		for _, s2 := range userInPartPendingGroups {
			if s2 != group.ID {
				newGroups = append(newGroups, s2)
			}
		}
		userInPartPendingGroupsJson, _ := json.Marshal(newGroups)
		initializers.DB.Model(&models.User{}).Where("id = ?", s).Update("pending_groups", userInPartPendingGroupsJson)
	}

	initializers.DB.Delete(&models.Group{}, group.ID)

	go realtime.GroupDeleted(group.GroupID, groupParticapants, groupPendingParticapants)

	c.JSON(http.StatusOK, gin.H{
		"message": "Success.",
	})
	return
}
