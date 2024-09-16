package models

import "gorm.io/gorm"

type Groups struct {
	gorm.Model
	GroupId     int
	GroupName   string
	GroupMember string
}
