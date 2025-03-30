/*
	Login with email
	This use case is used to login a user with email and password.
	The user will be logged in if the email and password are correct, otherwise an error will be returned.

	Steps:
		1. Validate the email and password
		2. Login with email and password
		3. Return the token
*/

package app

import (
	"context"
	"strings"

	securityDomain "github.com/TebanMT/smartGou/src/modules/security/domain"
)

type LoginWithEmail struct {
	securityRepository securityDomain.EmailAuthProvider
}

func NewLoginWithEmail(securityRepository securityDomain.EmailAuthProvider) *LoginWithEmail {
	return &LoginWithEmail{securityRepository: securityRepository}
}

func (l *LoginWithEmail) LoginWithEmail(ctx context.Context, email string, password string) (*securityDomain.TokenEntity, error) {
	email = strings.ToLower(email)
	token, err := l.securityRepository.LoginWithEmail(ctx, email, password)
	if err != nil {
		return nil, err
	}
	return token, nil
}
