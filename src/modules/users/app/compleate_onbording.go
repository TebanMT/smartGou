package app

import (
	"context"

	commonDomain "github.com/TebanMT/smartGou/src/common/domain"
	userDomain "github.com/TebanMT/smartGou/src/modules/users/domain"
	"github.com/google/uuid"
)

type CompleteOnboardingUseCase struct {
	userRepository userDomain.UserRepository
	unitOfWork     commonDomain.UnitOfWork
}

func NewCompleteOnboardingUseCase(userRepository userDomain.UserRepository, unitOfWork commonDomain.UnitOfWork) *CompleteOnboardingUseCase {
	return &CompleteOnboardingUseCase{userRepository: userRepository, unitOfWork: unitOfWork}
}

func (u *CompleteOnboardingUseCase) CompleteOnboarding(ctx context.Context, userID uuid.UUID) error {
	tx, err := u.unitOfWork.Begin(ctx)
	if err != nil {
		return err
	}

	rollback := true
	defer func() {
		if rollback {
			u.unitOfWork.Rollback(tx)
		}
	}()

	err = u.userRepository.CompleteOnboarding(tx, userID)
	if err != nil {
		return err
	}

	if err := u.unitOfWork.Commit(tx); err != nil {
		return err
	}
	rollback = false

	return nil
}
