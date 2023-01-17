package routes

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

type JobVolunteers struct {
	UserId    uint   `json:"userId"`
	Positions int8   `json:"positions"`
	FullName  string `json:"fullName"`
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
	var volunteers []uint
	volunteersJSON, _ := json.Marshal(volunteers)
	job := models.Job{Client: body.Client, Address: body.Address, Volunteer: string(volunteersJSON), Month: body.Date.Month, Day: body.Date.Day, Year: body.Date.Year, Hour: body.Time.Hour, Minute: body.Time.Minute, Pm: newPm, JobTitle: body.JobTitle, GroupId: group.GroupID, GroupNumId: group.ID, Instuctions: body.Instuctions, GroupName: group.Name, Issuer: user.ID, IssuerName: user.FirstName + " " + user.LastName, Taken: false, Positions: body.Positions}
	result := initializers.DB.Create(&job)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create job.",
		})
		return
	}

	if body.Notifications {
		realtime.NotifyPeople(group.ID, "Job Added", "Added a new Job")
		for _, s := range groupParticapants {
			var ownerSend models.User
			initializers.DB.First(&ownerSend, "id = ?", s)

			if len(ownerSend.NotificationToken) != 0 {
				// To check the token is valid
				pushToken, err := expo.NewExponentPushToken(ownerSend.NotificationToken)
				if err != nil {
					c.JSON(http.StatusOK, gin.H{
						"message": "Successfully Created The Job",
						"return":  job,
					})
					return
				}

				// Create a new Expo SDK client
				client := expo.NewPushClient(nil)

				// Publish message
				response, err := client.Publish(
					&expo.PushMessage{
						To:       []expo.ExponentPushToken{pushToken},
						Body:     user.FullName + " created a job for " + group.Name + ".",
						Data:     map[string]string{"groupId": group.GroupID, "type": "join"},
						Sound:    "default",
						Title:    "New job created in " + group.Name + ".",
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
		}
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

func GetJobsByMonthYear(c *gin.Context) {
	var body struct {
		Month int8  `json:"month"`
		Year  int16 `json:"year"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	if body.Month == 0 || body.Year == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing Parameters",
		})
		return
	}

	fmt.Println(body.Month, body.Year)

	user2, _ := c.Get("user")
	user := user2.(models.User)

	var userGroups []uint
	json.Unmarshal([]byte(user.Groups), &userGroups)

	fmt.Println(userGroups)

	var jobs []models.Job
	initializers.DB.Where("month = ? AND year = ? AND group_num_id IN ?", body.Month, body.Year, userGroups).Find(&jobs)

	c.JSON(http.StatusOK, gin.H{
		"jobs": jobs,
	})
	return
}

func JobInfo(c *gin.Context) {
	var body struct {
		JobId uint `json:"jobId"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	if body.JobId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing Parameters",
		})
		return
	}

	var job models.Job
	initializers.DB.First(&job, "id = ?", body.JobId)

	if job.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Group doesn't exist",
		})
		return
	}

	var group models.Group
	initializers.DB.First(&group, "group_id = ?", job.GroupId)

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

	c.JSON(http.StatusOK, job)
}

func AcceptJob(c *gin.Context) {
	var body struct {
		JobId     uint `json:"jobId"`
		Positions int8 `json:"positions"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	if body.JobId == 0 || body.Positions == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing Parameters",
		})
		return
	}

	var job models.Job
	initializers.DB.First(&job, "id = ?", body.JobId)

	if job.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Group doesn't exist",
		})
		return
	}

	var group models.Group
	initializers.DB.First(&group, "group_id = ?", job.GroupId)

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

	var jobVolunteers []JobVolunteers
	json.Unmarshal([]byte(job.Volunteer), &jobVolunteers)
	var positionsAlreadyTaken int16

	for _, s := range jobVolunteers {
		positionsAlreadyTaken = positionsAlreadyTaken + int16(s.Positions)
	}

	if positionsAlreadyTaken+int16(body.Positions) > int16(job.Positions) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"error": "Not enough positions.",
		})
		return
	}

	if positionsAlreadyTaken+int16(body.Positions) == int16(job.Positions) {
		initializers.DB.Model(&models.Job{}).Where("id = ?", job.ID).Update("taken", true)
	}

	var jobVolunteers2 []JobVolunteers
	alreadyAdded := false
	for _, s := range jobVolunteers {
		if s.UserId == user.ID {
			jobVolunteers2 = append(jobVolunteers2, JobVolunteers{FullName: s.FullName, Positions: s.Positions + body.Positions, UserId: s.UserId})
			alreadyAdded = true
		} else {
			jobVolunteers2 = append(jobVolunteers2, s)
		}
	}
	if alreadyAdded == false {
		jobVolunteers2 = append(jobVolunteers2, JobVolunteers{FullName: user.FullName, Positions: body.Positions, UserId: user.ID})
	}
	jobVolunteersJson, _ := json.Marshal(jobVolunteers2)
	initializers.DB.Model(&models.Job{}).Where("id = ?", job.ID).Update("volunteer", jobVolunteersJson)

	c.JSON(http.StatusOK, gin.H{
		"message": "Added new Volunteer.",
		"job":     job,
	})
}
