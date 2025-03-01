package models

import (
	"time"

	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	Id             int    `gorm:"primaryKey"`
	Username       string `gorm:"string;not null"`
	Name           string `gorm:"string;not null"`
	LastName       string `gorm:"string;not null"`
	SecondLastName string `gorm:"string;null"`
	Email          string `gorm:"string;not null"`
	DailingCode    string `gorm:"string;not null"`
	Phone          string `gorm:"string;not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
