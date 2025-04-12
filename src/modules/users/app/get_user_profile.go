/*
	GetUserProfile is a use case that gets the user profile for a user.
*/

package app

import (
	"context"

	"github.com/TebanMT/smartGou/src/modules/users/domain"
	commonDomain "github.com/TebanMT/smartGou/src/shared/domain"
	"github.com/google/uuid"
)

type GetUserProfile struct {
	UserRepository domain.UserRepository
	UnitOfWork     commonDomain.UnitOfWork
}

func NewGetUserProfile(userRepository domain.UserRepository, unitOfWork commonDomain.UnitOfWork) *GetUserProfile {
	return &GetUserProfile{
		UserRepository: userRepository,
		UnitOfWork:     unitOfWork,
	}
}

func (u *GetUserProfile) GetUserProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	tx, err := u.UnitOfWork.Query(ctx)
	if err != nil {
		return nil, err
	}

	user, err := u.UserRepository.GetUserByID(tx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
