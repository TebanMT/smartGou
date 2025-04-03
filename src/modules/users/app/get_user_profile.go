/*
	GetUserProfile is a use case that gets the user profile for a user.
*/

package app

import (
	"context"

	commonDomain "github.com/TebanMT/smartGou/src/common/domain"
	"github.com/TebanMT/smartGou/src/modules/users/domain"
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
	tx, err := u.UnitOfWork.Begin(ctx)
	if err != nil {
		return nil, err
	}

	rollback := true
	defer func() {
		if rollback {
			u.UnitOfWork.Rollback(tx)
		}
	}()

	user, err := u.UserRepository.GetUserByID(tx, userID)
	if err != nil {
		return nil, err
	}

	if err := u.UnitOfWork.Commit(tx); err != nil {
		return nil, err
	}
	rollback = false

	return user, nil
}
