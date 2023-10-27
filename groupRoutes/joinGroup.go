package grouproutes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	expo "github.com/oliveroneill/exponent-server-sdk-golang/sdk"
	"github.com/stephen10121/calenderapi/functions"
	"github.com/stephen10121/calenderapi/helpers"
	"github.com/stephen10121/calenderapi/initializers"
	"github.com/stephen10121/calenderapi/models"
	"github.com/stephen10121/calenderapi/realtime"
	"golang.org/x/crypto/bcrypt"
)

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

	userPendingGroups, err := helpers.UnmarshalPendingGroups(user.PendingGroups)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong",
		})
		fmt.Println(err)
		return
	}

	if functions.UintContains(userPendingGroups, group.ID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Already joined or attempted to join",
		})
		return
	}

	var userGroups []uint
	json.Unmarshal([]byte(user.Groups), &userGroups)
	if functions.UintContains(userGroups, group.ID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Already joined or attempted to join",
		})
		return
	}

	pendingGroups, err := helpers.UnmarshalPendingGroups(user.PendingGroups)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong",
		})
		fmt.Println(err)
		return
	}

	pendingGroups = append(pendingGroups, group.ID)
	pendingGroupsJson, _ := json.Marshal(pendingGroups)
	fmt.Println(pendingGroups)
	initializers.DB.Model(&models.User{}).Where("id = ?", user.ID).Update("pending_groups", pendingGroupsJson)

	var pendingParticapantSlice []uint
	json.Unmarshal([]byte(group.PendingParticapants), &pendingParticapantSlice)
	pendingParticapantSlice = append(pendingParticapantSlice, user.ID)
	pendingParticapantSliceJson, _ := json.Marshal(pendingParticapantSlice)
	initializers.DB.Model(&models.Group{}).Where("id = ?", group.ID).Update("pending_particapants", pendingParticapantSliceJson)

	var ownerSend models.User
	initializers.DB.First(&ownerSend, "id = ?", group.Owner)

	if len(ownerSend.NotificationToken) != 0 {
		// To check the token is valid
		pushToken, err := expo.NewExponentPushToken(ownerSend.NotificationToken)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"message":   "Success. Now wait for the group owner to accept the join request.",
				"groupName": group.Name,
			})
			return
		}

		// Create a new Expo SDK client
		client := expo.NewPushClient(nil)

		// Publish message
		response, err := client.Publish(
			&expo.PushMessage{
				To:       []expo.ExponentPushToken{pushToken},
				Body:     user.FullName + " wants to join your group.",
				Data:     map[string]string{"groupId": group.GroupID, "type": "join"},
				Sound:    "default",
				Title:    "New Join Request for " + group.Name,
				Priority: expo.HighPriority,
			},
		)

		// Check errors
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"message":   "Success. Now wait for the group owner to accept the join request.",
				"groupName": group.Name,
			})
			return
		}

		// Validate responses
		if response.ValidateResponse() != nil {
			fmt.Println(response.PushMessage.To, "failed")
		}
	}
	realtime.UserJoiningGroup(group.GroupID, user.FullName, group.Owner)

	c.JSON(http.StatusOK, gin.H{
		"message":   "Success. Now wait for the group owner to accept the join request.",
		"groupName": group.Name,
	})
	return
}
