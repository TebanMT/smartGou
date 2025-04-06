/*
Logout use case
Technically, this logout use case just invalidate the refresh token. The access token could be used until its expiration but you can´t create a new access token with the same refresh token.
Most of the entity properties does not have a direct logic to invalidate access token directly. We will need to implement a logic to invalidate access token but this
is not implemented yet. #TODO: To implement this logic we need to store some information in redis or some other database, and then we need to check if the information is alive.
Some ideas:
- Use blacklists to store the access tokens that have been invalidated.
- Use a TTL to invalidate the access token after a certain period of time.
- Use a logic to invalidate the access token when the user is logged out adding information in redis with the user id.
*/
package app

import (
	"context"

	securityDomain "github.com/TebanMT/smartGou/src/modules/security/domain"
)

type LogoutUseCase struct {
	securityRepository securityDomain.TokenManager
}

func NewLogoutUseCase(securityRepository securityDomain.TokenManager) *LogoutUseCase {
	return &LogoutUseCase{securityRepository: securityRepository}
}

func (u *LogoutUseCase) Logout(ctx context.Context, accessToken string) (bool, error) {
	if accessToken == "" {
		return false, securityDomain.ErrInvalidAccessToken
	}
	return u.securityRepository.Logout(ctx, accessToken)
}
