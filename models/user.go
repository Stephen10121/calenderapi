package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email             string `json:"email" gorm:"unique"`
	GoogId            string `json:"googId" gorm:"unique"`
	FirstName         string `json:"firstName"`
	LastName          string `json:"lastName"`
	FullName          string `json:"name"`
	Groups            string `json:"groups"`
	PendingGroups     string `json:"pendingGroups"`
	Locale            string `json:"locale"`
	Picture           string `json:"picture"`
	VerifiedEmail     bool   `json:"verifiedEmail"`
	NotificationToken string `json:"notificationToken"`
}
