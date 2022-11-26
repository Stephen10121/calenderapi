package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stephen10121/calenderapi/initializers"
	"github.com/stephen10121/calenderapi/models"
	"github.com/stephen10121/calenderapi/realtime"
)

func AddJob(c *gin.Context) {
	var body struct {
		Client        string `json:"client"`  //optional
		Address       string `json:"address"` //optional
		Date          string `json:"date"`
		Time          string `json:"time"`
		JobTitle      string `json:"jobTitle"`
		Group         string `json:"group"`
		Notifications bool   `json:"notifications"`
		Instuctions   string `json:"instructions"` //optional
		Positions     int8   `json:"positions"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	if body.Time == "" || body.Date == "" || body.JobTitle == "" || body.Group == "" || body.Positions == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing Parameters",
		})
		return
	}

	var group models.Group
	initializers.DB.First(&group, "group_id = ?", body.Group)

	if group.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Group doesn't exist",
		})
		return
	}

	user2, _ := c.Get("user")
	user := user2.(models.User)

	job := models.Job{Client: body.Client, Address: body.Address, Volunteer: "", Date: body.Date, Time: body.Time, JobTitle: body.JobTitle, Group: group.GroupID, Instuctions: body.Instuctions, GroupName: group.Name, Issuer: user.ID, IssuerName: user.Name, Taken: false}
	result := initializers.DB.Create(&job)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	if body.Notifications {
		realtime.NotifyPeople(group.ID, "Job Added", "Added a new Job")
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully Created The Job",
		"return":  job,
	})
	return
}
