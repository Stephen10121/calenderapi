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
	Group       string `json:"group"`
	Instuctions string `json:"instructions"`
	GroupName   string `json:"groupName"`
	Issuer      uint   `json:"issuer"`
	IssuerName  string `json:"issuerName"`
	Taken       bool   `json:"taken"`
}
