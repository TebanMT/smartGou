/*
	SignUpByEmail is a use case that signs up a user by email.
	The user will be created in the database and inderity provider (cognito) if it doesn't exist.
	If the user exists, an error will be returned.
	A mail will be sent to the user to verify the email.
	If the user exists but the email is not verified, a new OTP will be sent to the user.
*/

package app

import (
	"context"
	"strings"

	securityDomain "github.com/TebanMT/smartGou/src/modules/security/domain"
	userDomain "github.com/TebanMT/smartGou/src/modules/users/domain"
	commonDomain "github.com/TebanMT/smartGou/src/shared/domain"
	"github.com/TebanMT/smartGou/src/shared/utils"
)

type SignUpByEmailUseCase struct {
	securityRepository securityDomain.EmailAuthProvider
	userRepository     userDomain.UserRepository
	unitOfWork         commonDomain.UnitOfWork
}

func NewSignUpByEmailUseCase(securityRepository securityDomain.EmailAuthProvider, userRepository userDomain.UserRepository, unitOfWork commonDomain.UnitOfWork) *SignUpByEmailUseCase {
	return &SignUpByEmailUseCase{securityRepository: securityRepository, userRepository: userRepository, unitOfWork: unitOfWork}
}

func (u *SignUpByEmailUseCase) SignUpByEmail(ctx context.Context, email string, password string) (string, error) {

	plainPassword := password
	userEntity := userDomain.User{
		Email:         &email,
		Password:      &password,
		VerifiedPhone: false,
		VerifiedEmail: false,
	}
	err := utils.ValidateUserEmailAndPassword(&userEntity)
	if err != nil {
		return "", err
	}
	userEntity.HashPassword()

	tx, err := u.unitOfWork.Begin(ctx)
	if err != nil {
		return "", err
	}

	// Ensure rollback in case of error
	rollback := true
	defer func() {
		if rollback {
			u.unitOfWork.Rollback(tx)
		}
	}()

	user, err := utils.CheckUserExistenceByEmail(tx, u.userRepository, &userEntity)

	if err != nil && user == nil {
		return "", err
	}

	if user != nil && user.VerifiedEmail {
		return "", userDomain.ErrEmailAlreadyExists
	}

	if user == nil {
		userEntity.GenerateUserID()
		user, err = u.userRepository.CreateUser(tx, &userEntity)
		if err != nil {
			return "", err
		}
		email := strings.ToLower(*userEntity.Email)
		err := u.securityRepository.RegisterWithEmail(ctx, email, plainPassword, userEntity.UserID)
		if err != nil {
			return "", err
		}

		if err := u.unitOfWork.Commit(tx); err != nil {
			return "", err
		}

		rollback = false // ✅ DO NOT ROLLBACK IF EVERYTHING IS OK

		return user.UserID.String(), nil
	}

	err = u.securityRepository.ResendOtpByEmail(ctx, user.UserID)
	if err != nil {
		return "", err
	}

	return user.UserID.String(), nil
}
