package domain

import "errors"

var (
	ErrInvalidPhoneOrDailingCode = errors.New("invalid phone number or dailing code")
	ErrInvalidPhoneLength        = errors.New("phone number must be between 10 and 15 digits")
	ErrInvalidDailingCodeLength  = errors.New("dailing code must be between 1 and 5 digits")
	ErrInvalidDailingCodeFormat  = errors.New("dailing code must be a valid dailing code. Example: +52")
	ErrInvalidPhoneFormat        = errors.New("phone number must be a valid phone number. Example: 1234567890")
	ErrUserNotFoundException     = errors.New("user not found")
	ErrInvalidOTP                = errors.New("invalid OTP")
	ErrUserAlreadyConfirmed      = errors.New("user already confirmed")
	ErrUserAlreadyExists         = errors.New("user already exists")
	ErrExpiredOTP                = errors.New("OTP expired")
	ErrUserNotConfirmed          = errors.New("user not confirmed")
	ErrInvalidCredentials        = errors.New("invalid credentials")
	ErrRefreshTokenExpired       = errors.New("refresh token expired")
	ErrInvalidRefreshToken       = errors.New("invalid refresh token")
	ErrInvalidAccessToken        = errors.New("invalid access token")
	ErrInvalidSession            = errors.New("invalid session")
	ErrPhoneAlreadyVerified      = errors.New("phone already verified. Please login")
	ErrMaxAttemptsReached        = errors.New("max attempts reached")

	// Token errors
	ErrTokenExpired = errors.New("token expired")
	ErrInvalidToken = errors.New("invalid token")
	ErrUnauthorized = errors.New("unauthorized")
)
