package models

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	Owner        uint   `json:"owner"`
	OwnerEmail   string `json:"ownerEmail"`
	OwnerName    string `json:"ownerName"`
	GroupID      string `json:"groupID" gorm:"unique"`
	Password     string `json:"password"`
	Name         string `json:"name"`
	OthersCanAdd bool   `json:"othersCanAdd"`
	AboutGroup   string `json:"aboutGroup"`
}
