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

	group := models.Group{Owner: user.ID, OwnerEmail: user.Email, GroupID: body.Id, Password: string(hash), Name: body.Name, OthersCanAdd: body.OthersCanAdd, OwnerName: user.FirstName + " " + user.LastName, AboutGroup: body.AboutGroup, Particapants: strconv.FormatUint(uint64(user.ID), 10)}
	result := initializers.DB.Create(&group)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	groups := user.Groups + ":" + strconv.FormatUint(uint64(group.ID), 10)
	initializers.DB.Model(&models.User{}).Where("id = ?", user.ID).Update("groups", groups)

	users := strings.Split(group.Particapants, ":")
	var particapants []PartiacapantSend
	for _, s := range users {
		var user models.User
		initializers.DB.First(&user, "id = ?", s)
		if user.ID != 0 {
			particapants = append(particapants, PartiacapantSend{Name: user.FirstName + " " + user.LastName, Id: user.ID})
		}
	}

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

	users := strings.Split(group.Particapants, ":")
	if functions.Contains(users, strconv.FormatUint(uint64(user.ID), 10)) {
		var groupUsers []PartiacapantSend
		for _, s := range users {
			var user models.User
			initializers.DB.First(&user, "id = ?", s)
			if user.ID != 0 {
				groupUsers = append(groupUsers, PartiacapantSend{Name: user.FirstName + " " + user.LastName, Id: user.ID})
			}
		}

		if user.ID == group.Owner {
			usersPending := strings.Split(group.PendingParticapants, ":")
			var groupUsersPending []PartiacapantSend
			for _, s := range usersPending {
				var user models.User
				initializers.DB.First(&user, "id = ?", s)
				if user.ID != 0 {
					groupUsersPending = append(groupUsersPending, PartiacapantSend{Name: user.FirstName + " " + user.LastName, Id: user.ID})
				}
			}

			c.JSON(http.StatusOK, gin.H{
				"name":         group.Name,
				"owner":        group.OwnerName,
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
			"owner":        group.OwnerName,
			"created":      group.CreatedAt,
			"group_id":     group.GroupID,
			"about_group":  group.AboutGroup,
			"particapants": groupUsers,
			"yourowner":    false,
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

	messageToOwner := user.FirstName + " " + user.LastName + " wants to join your group."
	realtime.NotifyGroupOwner(group.ID, "Pending New User", messageToOwner)

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

	if user.PendingGroups != "" {
		usersCurrentPendingGroups := strings.Split(user.PendingGroups, ":")
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
	}

	var usersCurrentGroupsJson []groupData
	if user.Groups != "" {
		usersCurrentGroups := strings.Split(user.Groups, ":")

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
	}

	c.JSON(http.StatusOK, gin.H{
		"groups":        usersCurrentGroupsJson,
		"pendingGroups": userCurrentPendingGroupsJson,
	})
	return
}

func LeaveGroup(c *gin.Context) {
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

	groupUsers := strings.Split(group.Particapants, ":")
	if functions.Contains(groupUsers, strconv.FormatUint(uint64(user.ID), 10)) != true {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Your not part of the group.",
		})
		return
	}

	usersGroups := strings.Split(user.Groups, ":")
	var groups string

	if len(usersGroups) != 0 {
		for i := 0; i < len(usersGroups); i++ {
			if usersGroups[i] != strconv.FormatUint(uint64(group.ID), 10) && len(usersGroups[i]) != 0 {
				groups = groups + ":" + usersGroups[i]
			}
		}
	}
	initializers.DB.Model(&models.User{}).Where("id = ?", user.ID).Update("groups", groups)

	groupParticapants := strings.Split(group.Particapants, ":")
	var particapants string

	if len(groupParticapants) != 0 {
		for i := 0; i < len(groupParticapants); i++ {
			if groupParticapants[i] != strconv.FormatUint(uint64(user.ID), 10) && len(groupParticapants[i]) != 0 {
				particapants = particapants + ":" + groupParticapants[i]
			}
		}
	}
	initializers.DB.Model(&models.Group{}).Where("id = ?", group.ID).Update("particapants", particapants)

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

	groupUsers := strings.Split(group.PendingParticapants, ":")
	if functions.Contains(groupUsers, strconv.FormatUint(uint64(user.ID), 10)) != true {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Your not pending.",
		})
		return
	}

	usersGroups := strings.Split(user.PendingGroups, ":")
	var groups string

	if len(usersGroups) != 0 {
		for i := 0; i < len(usersGroups); i++ {
			if usersGroups[i] != strconv.FormatUint(uint64(group.ID), 10) && len(usersGroups[i]) != 0 {
				groups = groups + ":" + usersGroups[i]
			}
		}
	}
	initializers.DB.Model(&models.User{}).Where("id = ?", user.ID).Update("pending_groups", groups)

	groupParticapants := strings.Split(group.PendingParticapants, ":")
	var particapants string

	if len(groupParticapants) != 0 {
		for i := 0; i < len(groupParticapants); i++ {
			if groupParticapants[i] != strconv.FormatUint(uint64(user.ID), 10) && len(groupParticapants[i]) != 0 {
				particapants = particapants + ":" + groupParticapants[i]
			}
		}
	}
	initializers.DB.Model(&models.Group{}).Where("id = ?", group.ID).Update("pending_particapants", particapants)

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

	groupParticapants := strings.Split(group.Particapants, ":")
	if len(groupParticapants) != 0 {
		for i := 0; i < len(groupParticapants); i++ {
			u64, err := strconv.ParseUint(groupParticapants[i], 10, 16)
			if err != nil {
				fmt.Println(err)
				continue
			}

			var userInPart models.User
			initializers.DB.First(&user, "id = ?", u64)

			if userInPart.ID == 0 {
				continue
			}

			usersGroups := strings.Split(userInPart.Groups, ":")
			var newGroups string
			if len(usersGroups) != 0 {
				for i := 0; i < len(usersGroups); i++ {
					if usersGroups[i] != strconv.FormatUint(uint64(group.ID), 10) && len(usersGroups[i]) != 0 {
						newGroups = newGroups + ":" + usersGroups[i]
					}
				}
			}
			initializers.DB.Model(&models.User{}).Where("id = ?", u64).Update("groups", newGroups)
		}
	}

	groupPendingParticapants := strings.Split(group.PendingParticapants, ":")
	if len(groupPendingParticapants) != 0 {
		for i := 0; i < len(groupPendingParticapants); i++ {
			u64, err := strconv.ParseUint(groupPendingParticapants[i], 10, 16)
			if err != nil {
				fmt.Println(err)
				continue
			}

			var userInPart models.User
			initializers.DB.First(&user, "id = ?", u64)

			if userInPart.ID == 0 {
				continue
			}

			usersGroups := strings.Split(userInPart.PendingGroups, ":")
			var newGroups string
			if len(usersGroups) != 0 {
				for i := 0; i < len(usersGroups); i++ {
					if usersGroups[i] != strconv.FormatUint(uint64(group.ID), 10) && len(usersGroups[i]) != 0 {
						newGroups = newGroups + ":" + usersGroups[i]
					}
				}
			}
			initializers.DB.Model(&models.User{}).Where("id = ?", u64).Update("pending_groups", newGroups)
		}
	}

	usersGroups := strings.Split(user.PendingGroups, ":")
	var groups string

	if len(usersGroups) != 0 {
		for i := 0; i < len(usersGroups); i++ {
			if usersGroups[i] != strconv.FormatUint(uint64(group.ID), 10) && len(usersGroups[i]) != 0 {
				groups = groups + ":" + usersGroups[i]
			}
		}
	}
	initializers.DB.Model(&models.User{}).Where("id = ?", user.ID).Update("pending_groups", groups)

	var particapants string

	if len(groupParticapants) != 0 {
		for i := 0; i < len(groupParticapants); i++ {
			if groupParticapants[i] != strconv.FormatUint(uint64(user.ID), 10) && len(groupParticapants[i]) != 0 {
				particapants = particapants + ":" + groupParticapants[i]
			}
		}
	}

	initializers.DB.Delete(&models.Group{}, group.ID)

	realtime.GroupDeleted(group.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Success.",
	})
	return
}
