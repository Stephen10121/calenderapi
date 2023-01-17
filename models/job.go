package models

import (
	"gorm.io/gorm"
)

type Job struct {
	gorm.Model
	Client      string `json:"client"`
	Address     string `json:"address"`
	Volunteer   string `json:"volunteer"` // [{userId: 1, Positions: 2, FullName: Jeff Jeffins}]
	Month       int8   `json:"month"`
	Day         int8   `json:"day"`
	Year        int16  `json:"year"`
	Hour        int8   `json:"hour"`
	Minute      int8   `json:"minute"`
	Pm          bool   `json:"pm"`
	JobTitle    string `json:"jobTitle"`
	GroupId     string `json:"groupId"`
	GroupNumId  uint   `json:"groupNumId"`
	Instuctions string `json:"instructions"`
	GroupName   string `json:"groupName"`
	Issuer      uint   `json:"issuer"`
	IssuerName  string `json:"issuerName"`
	Taken       bool   `json:"taken"`
	Positions   int8   `json:"positions"`
}
