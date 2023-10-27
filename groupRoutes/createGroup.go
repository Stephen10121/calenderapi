package grouproutes

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stephen10121/calenderapi/initializers"
	"github.com/stephen10121/calenderapi/models"
	"golang.org/x/crypto/bcrypt"
)

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
