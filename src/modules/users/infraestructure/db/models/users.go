package models

import (
	"time"

	"github.com/google/uuid"
)

type UserModel struct {
	Id                    int       `gorm:"primaryKey"`
	UserID                uuid.UUID `gorm:"type:uuid;not null;unique"`
	Username              *string   `gorm:"string;null"`
	Name                  *string   `gorm:"string;null"`
	LastName              *string   `gorm:"string;null"`
	SecondLastName        *string   `gorm:"string;null"`
	Email                 *string   `gorm:"string;null;unique; default:null"`
	DailingCode           *string   `gorm:"string;null"`
	Phone                 *string   `gorm:"string;null;unique"`
	Password              *string   `gorm:"-"`
	IsOnboardingCompleted bool      `gorm:"boolean;default:false"`
	VerifiedPhone         bool      `gorm:"boolean;default:false"`
	VerifiedEmail         bool      `gorm:"boolean;default:false"`
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

func (UserModel) TableName() string {
	return "users"
}
