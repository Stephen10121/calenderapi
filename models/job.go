package models

import (
	"gorm.io/gorm"
)

type Job struct {
	gorm.Model
	Client      string `json:"client"`
	Address     string `json:"address"`
	Volunteer   string `json:"volunteer"`
	Date        string `json:"date"`
	Time        string `json:"time"`
	JobTitle    string `json:"jobTitle"`
	Group       uint   `json:"group"`
	Instuctions string `json:"instructions"`
	GroupName   string `json:"groupName"`
	Issuer      string `json:"issuer"`
	Taken       bool   `json:"taken"`
}
