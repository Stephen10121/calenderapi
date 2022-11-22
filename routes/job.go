package routes

import "github.com/gin-gonic/gin"

type date struct {
	Year  int16 `json:"year"`
	Month int8  `json:"month"`
	Day   int8  `json:"day"`
}

type timeType struct {
	Hour   int8 `json:"hour"`
	Minute int8 `json:"minute"`
	PM     bool `json:"pm"`
}

func AddJob(c *gin.Context) {
	var body struct {
		Client      string   `json:"client"`
		Address     string   `json:"address"`
		Date        date     `json:"date"`
		Time        timeType `json:"time"`
		JobTitle    string   `json:"jobTitle"`
		Group       uint     `json:"group"`
		Instuctions string   `json:"instructions"`
	}
}
