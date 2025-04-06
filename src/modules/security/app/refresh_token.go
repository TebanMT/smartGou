/*
RefreshTokenUseCase is a use case that refreshes a token.
The access token is expired, so we need to refresh it.
*/
package app

import (
	"context"

	securityDomain "github.com/TebanMT/smartGou/src/modules/security/domain"
)

type RefreshTokenUseCase struct {
	tokenService securityDomain.TokenManager
}

func NewRefreshTokenUseCase(tokenService securityDomain.TokenManager) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{tokenService: tokenService}
}

func (u *RefreshTokenUseCase) RefreshToken(ctx context.Context, refreshToken string) (*securityDomain.TokenEntity, error) {
	token, err := u.tokenService.RefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	return token, nil
}
