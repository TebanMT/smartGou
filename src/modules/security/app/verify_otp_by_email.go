/*
VerifyOTPByEmail is a use case that verifies a OTP by email.
Verify the code given a userID and a code. The code is the code that the user received in the email.
Also, the email will be verified as true in provider (cognito) and database.
*/
package app

import (
	"context"

	securityDomain "github.com/TebanMT/smartGou/src/modules/security/domain"
	userDomain "github.com/TebanMT/smartGou/src/modules/users/domain"
	commonDomain "github.com/TebanMT/smartGou/src/shared/domain"
	"github.com/google/uuid"
)

type VerifyOTPByEmail struct {
	securityRepository securityDomain.EmailAuthenticator
	userRepository     userDomain.UserRepository
	unitOfWork         commonDomain.UnitOfWork
}

func NewVerifyOTPByEmail(securityRepository securityDomain.EmailAuthenticator, userRepository userDomain.UserRepository, unitOfWork commonDomain.UnitOfWork) *VerifyOTPByEmail {
	return &VerifyOTPByEmail{securityRepository: securityRepository, userRepository: userRepository, unitOfWork: unitOfWork}
}

func (c *VerifyOTPByEmail) VerifyOTPByEmail(ctx context.Context, userID uuid.UUID, code string) error {
	tx, err := c.unitOfWork.Begin(ctx)
	if err != nil {
		return err
	}

	rollback := true
	defer func() {
		if rollback {
			c.unitOfWork.Rollback(tx)
		}
	}()

	err = c.userRepository.VerifyEmail(tx, userID)
	if err != nil {
		return err
	}

	err = c.securityRepository.ConfirmOtpByEmail(ctx, userID, code)
	if err != nil {
		return err
	}

	c.unitOfWork.Commit(tx)
	rollback = false

	return nil
}
