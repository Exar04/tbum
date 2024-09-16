package models

import (
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique; column:username" json:"username"`
	// Email         string `gorm:"unique"`
	Password      string `json:"password"`
	Phone         string
	Token         string
	User_type     string
	Refresh_token string
}

type Claims struct {
	Username string
	jwt.StandardClaims
}
