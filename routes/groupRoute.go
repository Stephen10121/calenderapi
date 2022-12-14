package routes

// u64, err := strconv.ParseUint(s, 10, 16)
// if err != nil {
// 	continue
// }

//strconv.FormatUint(uint64(group.ID), 10)

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	expo "github.com/oliveroneill/exponent-server-sdk-golang/sdk"
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
		Name         string `json:"name"`
		Id           string `json:"id"`
		Password     string `json:"password"`
		OthersCanAdd bool   `json:"othersCanAdd"`
		AboutGroup   string `json:"aboutGroup"`
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
	particapantsJson, _ := json.Marshal([]uint{user.ID})
	pendingPartJson, _ := json.Marshal([]uint{})
	group := models.Group{Owner: user.ID, OwnerEmail: user.Email, GroupID: body.Id, Password: string(hash), Name: body.Name, OthersCanAdd: body.OthersCanAdd, OwnerName: user.FirstName + " " + user.LastName, AboutGroup: body.AboutGroup, Particapants: string(particapantsJson), PendingParticapants: string(pendingPartJson)}
	result := initializers.DB.Create(&group)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	var groups []uint
	json.Unmarshal([]byte(user.Groups), &groups)
	groups = append(groups, group.ID)
	groupsJson, _ := json.Marshal(groups)
	initializers.DB.Model(&models.User{}).Where("id = ?", user.ID).Update("groups", groupsJson)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"groupName":    group.Name,
			"groupOwner":   group.OwnerName,
			"groupId":      group.GroupID,
			"othersCanAdd": group.OthersCanAdd,
		},
	})
	return
}

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
	var groupParticapants []uint
	json.Unmarshal([]byte(group.Particapants), &groupParticapants)
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
	var userPendingGroups []uint
	json.Unmarshal([]byte(user.PendingGroups), &userPendingGroups)
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

	var pendingGroups []uint
	json.Unmarshal([]byte(user.PendingGroups), &pendingGroups)
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
				Data:     map[string]string{"withSome": "data"},
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

	var userCurrentPendingGroupsJson []groupPendingData
	var userPendingGroups []uint
	json.Unmarshal([]byte(user.PendingGroups), &userPendingGroups)
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
		usersCurrentGroupsJson = append(usersCurrentGroupsJson, groupData{GroupId: group.GroupID, GroupName: group.Name, OthersCanAdd: group.OthersCanAdd, GroupOwner: group.OwnerName})
	}

	c.JSON(http.StatusOK, gin.H{
		"groups":        usersCurrentGroupsJson,
		"pendingGroups": userCurrentPendingGroupsJson,
	})
	return
}

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

		var groupParticapants []uint
		json.Unmarshal([]byte(group.Particapants), &groupParticapants)
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
	var groupParticapants []uint
	json.Unmarshal([]byte(group.Particapants), &groupParticapants)
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
	var userPendingGroups []uint
	json.Unmarshal([]byte(user.PendingGroups), &userPendingGroups)
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
		var userInPartPendingGroups []uint
		json.Unmarshal([]byte(userInPart.PendingGroups), &userInPartPendingGroups)
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
