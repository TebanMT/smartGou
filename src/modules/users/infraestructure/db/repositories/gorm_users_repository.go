package repositories

import (
	"github.com/TebanMT/smartGou/src/modules/users/domain"
	"github.com/TebanMT/smartGou/src/modules/users/infraestructure/db/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *domain.User) error {
	userModel := models.UserModel{
		Id:             user.ID,
		Username:       user.Username,
		Name:           user.Name,
		LastName:       user.LastName,
		SecondLastName: user.SecondLastName,
		Email:          user.Email,
		DailingCode:    user.DailingCode,
		Phone:          user.Phone,
	}
	return r.db.Create(&userModel).Error
}
