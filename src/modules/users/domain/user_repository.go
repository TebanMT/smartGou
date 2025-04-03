package domain

import (
	"github.com/TebanMT/smartGou/src/common/domain"
	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(tx domain.Transaction, user *User) (*User, error)
	DeleteUser(tx domain.Transaction, userID uuid.UUID) error
	ExistsUserByEmail(tx domain.Transaction, email string) (*User, error)
	ExistsUserByPhone(tx domain.Transaction, phone string) (*User, error)
	UpdateUser(tx domain.Transaction, user *User) (*User, error)
	CompleteOnboarding(tx domain.Transaction, userID uuid.UUID) error
	VerifyPhone(tx domain.Transaction, userID uuid.UUID) error
	VerifyEmail(tx domain.Transaction, userID uuid.UUID) error
	GetUserByID(tx domain.Transaction, userID uuid.UUID) (*User, error)
}
