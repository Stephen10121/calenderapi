package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `json:"username" gorm:"unique"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Groups   string `json:"groups"`
}
