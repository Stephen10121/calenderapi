package models

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	Owner        uint   `json:"owner"`
	OwnerEmail   string `json:"ownerEmail"`
	GroupID      string `json:"groupID"`
	Password     string `json:"password"`
	Name         string `json:"name"`
	OthersCanAdd bool   `json:"othersCanAdd"`
}
