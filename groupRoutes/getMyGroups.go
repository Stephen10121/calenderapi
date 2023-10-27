package grouproutes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stephen10121/calenderapi/helpers"
	"github.com/stephen10121/calenderapi/initializers"
	"github.com/stephen10121/calenderapi/models"
)

type groupPendingData struct {
	GroupId   string `json:"groupId"`
	GroupName string `json:"groupName"`
}

type groupData struct {
	GroupId      string `json:"groupId"`
	GroupName    string `json:"groupName"`
	GroupOwner   string `json:"groupOwner"`
	YouOwn       bool   `json:"youOwn"`
	OthersCanAdd bool   `json:"othersCanAdd"`
}

func GetMyGroups(c *gin.Context) {
	user2, _ := c.Get("user")
	user := user2.(models.User)

	var userCurrentPendingGroupsJson []groupPendingData

	userPendingGroups, err := helpers.UnmarshalPendingGroups(user.PendingGroups)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong",
		})
		fmt.Println(err)
		return
	}

	for _, s := range userPendingGroups {
		var group models.Group
		initializers.DB.First(&group, s)

		if group.ID == 0 {
			continue
		}
		userCurrentPendingGroupsJson = append(userCurrentPendingGroupsJson, groupPendingData{GroupId: group.GroupID, GroupName: group.Name})
	}

	var usersCurrentGroupsJson []groupData
	var userGroups []uint
	json.Unmarshal([]byte(user.Groups), &userGroups)
	for _, s := range userGroups {
		var group models.Group
		initializers.DB.First(&group, s)

		if group.ID == 0 {
			continue
		}
		var iown = false
		if group.Owner == user.ID {
			iown = true
		}
		usersCurrentGroupsJson = append(usersCurrentGroupsJson, groupData{GroupId: group.GroupID, GroupName: group.Name, OthersCanAdd: group.OthersCanAdd, GroupOwner: group.OwnerName, YouOwn: iown})
	}

	c.JSON(http.StatusOK, gin.H{
		"groups":        usersCurrentGroupsJson,
		"pendingGroups": userCurrentPendingGroupsJson,
	})
	return
}
