package domain

import (
	"time"

	"github.com/google/uuid"
)

type TokenEntity struct {
	AccessToken  string
	RefreshToken string
	IdToken      string
	ExpiresIn    int
}

type LoginChallengeByPhone struct {
	Code        string // OTP code
	Session     string
	UserId      uuid.UUID
	MaxAttempts int
	ExpiresAt   time.Time
}

type VerifyOTPByPhoneResponse struct {
	TokenEntity    *TokenEntity
	LoginChallenge *LoginChallengeByPhone
}

func NewLoginChallengeByPhone(code string, session string, userId uuid.UUID) *LoginChallengeByPhone {
	return &LoginChallengeByPhone{
		Code:        code,
		Session:     session,
		UserId:      userId,
		MaxAttempts: 3,
		ExpiresAt:   time.Now().Add(time.Minute * 2),
	}
}

func NewTokenEntity(accessToken string, refreshToken string, idToken string, expiresIn int) *TokenEntity {
	return &TokenEntity{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		IdToken:      idToken,
		ExpiresIn:    expiresIn,
	}
}

func (t *TokenEntity) GetAccessToken() string {
	return t.AccessToken
}

func (t *TokenEntity) GetRefreshToken() string {
	return t.RefreshToken
}

func (t *TokenEntity) GetIdToken() string {
	return t.IdToken
}

func (t *TokenEntity) GetExpiresIn() int {
	return t.ExpiresIn
}
