/*
	RequestOTPByPhone is a use case that requests a OTP by phone number.
	The user will be created in the database and inderity provider (cognito) if it doesn't exist.
	If the user exists, the code will be sent to the phone number.
	The user will be created without password and return the session.
	The session is specific context for the user to login in the cognito user pool. It is used to verify the user's phone number.
	For others implementations, this session could be not needed.
	The user ID (uuid type) is the user's ID in the database and also the username in the cognito user pool.
*/

package app

import (
	"context"

	"github.com/TebanMT/smartGou/src/common"
	commonDomain "github.com/TebanMT/smartGou/src/common/domain"
	securityDomain "github.com/TebanMT/smartGou/src/modules/security/domain"
	userDomain "github.com/TebanMT/smartGou/src/modules/users/domain"
)

type RequestOTPByPhone struct {
	securityRepository securityDomain.PhoneAuthProvider
	userRepository     userDomain.UserRepository
	unitOfWork         commonDomain.UnitOfWork
}

func NewRequestOTPByPhone(securityRepository securityDomain.PhoneAuthProvider, userRepository userDomain.UserRepository, unitOfWork commonDomain.UnitOfWork) *RequestOTPByPhone {
	return &RequestOTPByPhone{securityRepository: securityRepository, userRepository: userRepository, unitOfWork: unitOfWork}
}

func (r *RequestOTPByPhone) RequestOTPByPhone(ctx context.Context, phoneNumber string, dailingCode string) (*securityDomain.LoginChallengeByPhone, error) {
	securityEntity := securityDomain.SecurityEntity{
		PhoneNumber: phoneNumber,
		DailingCode: dailingCode,
	}

	userEntity := userDomain.User{
		Phone:         &phoneNumber,
		DailingCode:   &dailingCode,
		VerifiedPhone: false,
		VerifiedEmail: false,
	}

	err := securityEntity.ValidatePhoneNumber()
	if err != nil {
		return nil, err
	}

	formattedPhoneNumber := securityEntity.FormatPhoneE164()

	tx, err := r.unitOfWork.Begin(ctx)
	if err != nil {
		return nil, err
	}

	// Ensure rollback in case of error
	rollback := true
	defer func() {
		if rollback {
			r.unitOfWork.Rollback(tx)
		}
	}()

	user, err := common.CheckUserExistenceByPhone(tx, r.userRepository, &userEntity)

	if err != nil && user == nil {
		return nil, err
	}

	if user.ID == 0 {
		userEntity.GenerateUserID()
		user, err = r.userRepository.CreateUser(tx, &userEntity)
		if err != nil {
			return nil, err
		}
		err := r.securityRepository.RegisterWithPhoneNumber(ctx, formattedPhoneNumber, user.UserID)
		if err != nil {
			return nil, err
		}
	}

	loginChallenge, err := r.securityRepository.SendOTPToPhone(ctx, user.UserID)
	if err != nil {
		return nil, err
	}

	if err := r.unitOfWork.Commit(tx); err != nil {
		return nil, err
	}

	rollback = false // ✅ DO NOT ROLLBACK IF EVERYTHING IS OK

	return loginChallenge, nil
}
