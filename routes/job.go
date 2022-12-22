package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stephen10121/calenderapi/functions"
	"github.com/stephen10121/calenderapi/initializers"
	"github.com/stephen10121/calenderapi/models"
	"github.com/stephen10121/calenderapi/realtime"
)

type TimeType struct {
	Hour   int8   `json:"hour"`
	Minute int8   `json:"minute"`
	Pm     string `json:"pm"`
}

type DateType struct {
	Month int8  `json:"month"`
	Day   int8  `json:"day"`
	Year  int16 `json:"year"`
}

func AddJob(c *gin.Context) {
	var body struct {
		Client        string   `json:"client"`  //optional
		Address       string   `json:"address"` //optional
		Date          DateType `json:"date"`
		Time          TimeType `json:"time"`
		JobTitle      string   `json:"jobTitle"`
		Group         string   `json:"group"`
		Notifications bool     `json:"notifications"`
		Instuctions   string   `json:"instructions"` //optional
		Positions     int8     `json:"positions"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	if body.JobTitle == "" || body.Group == "" || body.Positions == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing Parameters",
		})
		return
	}

	if body.Date.Month == 0 || body.Date.Day == 0 || body.Date.Year == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "The Date is invalid.",
		})
		return
	}

	if body.Time.Hour == 0 || body.Time.Minute == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "The Time is invalid.",
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

	var groupParticapants []uint
	json.Unmarshal([]byte(group.Particapants), &groupParticapants)

	if user.ID != group.Owner {
		if functions.UintContains(groupParticapants, user.ID) != true || group.OthersCanAdd != true {
			c.JSON(http.StatusMethodNotAllowed, gin.H{
				"error": "User not allowed to add job",
			})
			return
		}
	}

	var newPm = false
	if body.Time.Pm == "PM" {
		newPm = true
	}

	job := models.Job{Client: body.Client, Address: body.Address, Volunteer: "", Month: body.Date.Month, Day: body.Date.Day, Year: body.Date.Year, Hour: body.Time.Hour, Minute: body.Time.Minute, Pm: newPm, JobTitle: body.JobTitle, GroupId: group.GroupID, Instuctions: body.Instuctions, GroupName: group.Name, Issuer: user.ID, IssuerName: user.FirstName + " " + user.LastName, Taken: false, Positions: body.Positions}
	result := initializers.DB.Create(&job)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create job.",
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

func GetJobs(c *gin.Context) {
	var body struct {
		Group string `json:"group"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	if len(body.Group) == 0 {
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

	var groupParticapants []uint
	json.Unmarshal([]byte(group.Particapants), &groupParticapants)
	if functions.UintContains(groupParticapants, user.ID) != true {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"error": "User not part of group",
		})
		return
	}

	var jobs []models.Job
	initializers.DB.Where("group_id = ?", body.Group).Find(&jobs)

	c.JSON(http.StatusOK, gin.H{
		"jobs": jobs,
	})
	return
}
