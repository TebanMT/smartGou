package repositories

import (
	"fmt"
	"strings"

	"github.com/TebanMT/smartGou/src/common"
	commonDomain "github.com/TebanMT/smartGou/src/common/domain"
	"github.com/TebanMT/smartGou/src/modules/users/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) CreateUser(tx commonDomain.Transaction, user *domain.User) (*domain.User, error) {
	fmt.Println("Creating user", user)
	err := tx.Execute(func(t commonDomain.Transaction) error {
		gormTx := t.(*common.GormTransaction)
		err := gormTx.Tx.Create(&user).Error
		if err != nil {
			return err
		}

		return nil
	})

	fmt.Println("Error creating user", user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) ExistsUserByEmail(tx commonDomain.Transaction, email string) (*domain.User, error) {
	var user domain.User
	err := tx.Execute(func(t commonDomain.Transaction) error {
		gormTx := t.(*common.GormTransaction)
		return gormTx.Tx.Model(&domain.User{}).Where("LOWER(email) = ?", strings.ToLower(email)).First(&user).Error
	})
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) ExistsUserByPhone(tx commonDomain.Transaction, phone string) (*domain.User, error) {
	var user domain.User
	err := tx.Execute(func(t commonDomain.Transaction) error {
		gormTx := t.(*common.GormTransaction)
		return gormTx.Tx.Model(&domain.User{}).Where("phone = ?", phone).First(&user).Error
	})
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) DeleteUser(tx commonDomain.Transaction, userID uuid.UUID) error {
	return tx.Execute(func(t commonDomain.Transaction) error {
		gormTx := t.(*common.GormTransaction)
		return gormTx.Tx.Delete(&domain.User{}, "user_id = ?", userID).Error
	})
}

func (r *UserRepository) UpdateUser(tx commonDomain.Transaction, user *domain.User) (*domain.User, error) {
	err := tx.Execute(func(t commonDomain.Transaction) error {
		gormTx := t.(*common.GormTransaction)
		return gormTx.Tx.Model(&domain.User{}).Where("user_id = ?", user.UserID).Omit("id", "created_at", "updated_at").Updates(user).Error
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) CompleteOnboarding(tx commonDomain.Transaction, userID uuid.UUID) error {
	return tx.Execute(func(t commonDomain.Transaction) error {
		gormTx := t.(*common.GormTransaction)
		return gormTx.Tx.Model(&domain.User{}).Where("user_id = ?", userID).Update("is_onboarding_completed", true).Error
	})
}

func (r *UserRepository) VerifyPhone(tx commonDomain.Transaction, userID uuid.UUID) error {
	return tx.Execute(func(t commonDomain.Transaction) error {
		gormTx := t.(*common.GormTransaction)
		return gormTx.Tx.Model(&domain.User{}).Where("user_id = ?", userID).Update("verified_phone", true).Error
	})
}

func (r *UserRepository) VerifyEmail(tx commonDomain.Transaction, userID uuid.UUID) error {
	return tx.Execute(func(t commonDomain.Transaction) error {
		gormTx := t.(*common.GormTransaction)
		return gormTx.Tx.Model(&domain.User{}).Where("user_id = ?", userID).Update("verified_email", true).Error
	})
}
