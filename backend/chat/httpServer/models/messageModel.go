package models

import "gorm.io/gorm"

type Messages struct {
	gorm.Model
	Sender   string
	Reciever string
	Content  string
}
