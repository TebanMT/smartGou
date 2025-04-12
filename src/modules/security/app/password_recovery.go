/*
RequestPasswordRecovery
*/
package app

import (
	"context"
	"fmt"
	"strings"

	securityDomain "github.com/TebanMT/smartGou/src/modules/security/domain"
	userDomain "github.com/TebanMT/smartGou/src/modules/users/domain"
	commonDomain "github.com/TebanMT/smartGou/src/shared/domain"
	"github.com/TebanMT/smartGou/src/shared/utils"
)

type passwordRecoveryUseCase struct {
	passwordRecoveryRepository securityDomain.PasswordManager
	userRepository             userDomain.UserRepository
	unitOfWork                 commonDomain.UnitOfWork
}

func NewPasswordRecoveryUseCase(passwordRecoveryRepository securityDomain.PasswordManager, userRepository userDomain.UserRepository, unitOfWork commonDomain.UnitOfWork) *passwordRecoveryUseCase {
	return &passwordRecoveryUseCase{
		passwordRecoveryRepository: passwordRecoveryRepository,
		userRepository:             userRepository,
		unitOfWork:                 unitOfWork,
	}
}

func (u *passwordRecoveryUseCase) RequestPasswordRecovery(ctx context.Context, email string) (bool, error) {
	userValidator := userDomain.BuildUserValidatorEmailChain()
	user, err := u.ValidateUserByEmail(ctx, email, "", userValidator)
	if err != nil {
		return false, err
	}
	return u.passwordRecoveryRepository.PasswordRecovery(ctx, strings.ToLower(*user.Email))
}

func (u *passwordRecoveryUseCase) ResetPassword(ctx context.Context, email string, newPassword string, confirmationCode string) (bool, error) {
	userValidator := userDomain.BuildUserValidatorEmailAndPasswordChain()
	user, err := u.ValidateUserByEmail(ctx, email, newPassword, userValidator)
	if err != nil {
		return false, err
	}
	tx, err := u.unitOfWork.Begin(ctx)
	if err != nil {
		return false, err
	}

	rollback := true
	defer func() {
		if rollback {
			u.unitOfWork.Rollback(tx)
		}
	}()

	plainPassword := newPassword
	user.Password = &newPassword
	user.HashPassword()

	user, err = u.userRepository.UpdateUser(tx, user)
	if err != nil {
		return false, err
	}

	u.unitOfWork.Commit(tx)
	rollback = false
	return u.passwordRecoveryRepository.PasswordReset(ctx, user.UserID, plainPassword, confirmationCode)
}

func (u *passwordRecoveryUseCase) ValidateUserByEmail(ctx context.Context, email string, password string, validator userDomain.UserValidator) (*userDomain.User, error) {
	if email == "" {
		return nil, userDomain.ErrInvalidEmail
	}
	user := &userDomain.User{Email: &email, Password: &password}
	err := validator.Validate(user)
	if err != nil {
		return nil, err
	}
	tx, err := u.unitOfWork.Begin(ctx)
	if err != nil {
		return nil, err
	}

	rollback := true
	defer func() {
		if rollback {
			u.unitOfWork.Rollback(tx)
		}
	}()
	user, err = utils.CheckUserExistenceByEmail(tx, u.userRepository, user)
	if user == nil {
		fmt.Println("User not found. Error: ", err)
		return nil, userDomain.ErrUserNotFound
	}

	u.unitOfWork.Commit(tx)
	rollback = false

	return user, nil
}
