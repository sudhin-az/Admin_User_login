package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FullName    string `gorm:"not null"`
	UserName    string `gorm:"unique;not null"`
	Email       string `gorm:"not null"`
	PhoneNumber string `gorm:"not null"`
	Password    string `gorm:"not null"`
	Gender      string `gorm:"not null"`
	Role        string `gorm:"default:'user'"`
}

type UserDetails struct {
	UserName string
	Email    string
}
