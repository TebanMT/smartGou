/*
	VerifyOTPByPhone is a use case that verifies a OTP by phone number.
	If the OTP is valid, the user will be logged in and return the token entity.
	If the OTP is invalid, the user will be logged in and return the login challenge entity.
	The login challenge entity contains the session and the userID with which the user can try to verify the OTP again.
*/

package app

import (
	"context"

	commonDomain "github.com/TebanMT/smartGou/src/common/domain"
	securityDomain "github.com/TebanMT/smartGou/src/modules/security/domain"
	userDomain "github.com/TebanMT/smartGou/src/modules/users/domain"
	"github.com/google/uuid"
)

type VerifyOTPByPhone struct {
	securityRepository securityDomain.PhoneAuthProvider
	userRepository     userDomain.UserRepository
	unitOfWork         commonDomain.UnitOfWork
}

func NewVerifyOTPByPhone(securityRepository securityDomain.PhoneAuthProvider, userRepository userDomain.UserRepository, unitOfWork commonDomain.UnitOfWork) *VerifyOTPByPhone {
	return &VerifyOTPByPhone{securityRepository: securityRepository, userRepository: userRepository, unitOfWork: unitOfWork}
}

func (v *VerifyOTPByPhone) VerifyOTPByPhone(ctx context.Context, code string, session string, userID uuid.UUID) (*securityDomain.TokenEntity, *securityDomain.LoginChallengeByPhone, error) {
	loginChallenge := securityDomain.NewLoginChallengeByPhone(code, session, userID)

	tx, err := v.unitOfWork.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Ensure rollback in case of error
	rollback := true
	defer func() {
		if rollback {
			v.unitOfWork.Rollback(tx)
		}
	}()

	verifyOTPByPhoneResponse, err := v.securityRepository.VerifyOTPFromPhone(ctx, loginChallenge)

	if err != nil && verifyOTPByPhoneResponse.LoginChallenge != nil {
		return nil, verifyOTPByPhoneResponse.LoginChallenge, err
	}

	err = v.userRepository.VerifyPhone(tx, userID)

	if err != nil {
		return nil, nil, err
	}

	if err := v.unitOfWork.Commit(tx); err != nil {
		return nil, nil, err
	}

	rollback = false // ✅ DO NOT ROLLBACK IF EVERYTHING IS OK

	return verifyOTPByPhoneResponse.TokenEntity, nil, nil
}
