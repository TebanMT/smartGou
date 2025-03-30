package domain

import (
	"context"

	"github.com/google/uuid"
)

type PhoneAuthenticator interface {
	/*
		Send a code to the user's phone number given a userID. The user will be created without password and return the session
		The userID is the user's ID in the database and also the username in the cognito user pool
		The session is specific context for the user to login in the cognito user pool. It is used to verify the user's phone number.
		For others implementations, this session could be not needed.
	*/
	SendOTPToPhone(ctx context.Context, userID uuid.UUID) (*LoginChallengeByPhone, error)
	/*
		Verify the code given a userID, a code and a session. The code is the code that the user received in the phone number
		The session is the session that the user received in the SendOTPToPhone method. Again, for others implementations, this session could be not needed.
	*/
	VerifyOTPFromPhone(ctx context.Context, loginChallenge *LoginChallengeByPhone) (*VerifyOTPByPhoneResponse, error)
}

type EmailAuthenticator interface {
	/*
		Login with email and password
	*/
	LoginWithEmail(ctx context.Context, email string, password string) (*TokenEntity, error)
	/*
		Confirm the OTP by email
		Verify the code given a userID and a code. The code is the code that the user received in the email.
		Also, the email will be verified as true in provider (cognito) and database.
	*/
	ConfirmOtpByEmail(ctx context.Context, userID uuid.UUID, code string) error
	/*
		Resend the OTP by email
		Resend the OTP by email given a userID. The OTP will be sent to the user's email.
	*/
	ResendOtpByEmail(ctx context.Context, userID uuid.UUID) error
}

type IdentityRegister interface {
	/*
		Register with phone number
	*/
	RegisterWithPhoneNumber(ctx context.Context, phoneNumber string, userID uuid.UUID) error
	/*
		Register with email
	*/
	RegisterWithEmail(ctx context.Context, email string, password string, userID uuid.UUID) error
}

type TokenManager interface {
	/*
		Refresh the token
	*/
	RefreshToken(ctx context.Context, refreshToken string) (*TokenEntity, error)
	/*
		Logout the user
	*/
	Logout(ctx context.Context, accessToken string) (bool, error)
}

type PhoneAuthProvider interface {
	PhoneAuthenticator
	IdentityRegister
}

type EmailAuthProvider interface {
	EmailAuthenticator
	IdentityRegister
}
