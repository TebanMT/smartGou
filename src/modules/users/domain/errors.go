package domain

import "errors"

var (
	ErrPhoneOrEmailAlreadyExists = errors.New("User already exists with this phone or email")
	ErrEmailAlreadyExists        = errors.New("User already exists with this email")
	ErrPhoneAlreadyExists        = errors.New("User already exists with this phone")
	ErrUsernameAlreadyExists     = errors.New("User already exists with this username")
	ErrEmailRequired             = errors.New("email is required")
	ErrPhoneRequired             = errors.New("phone is required")
	ErrUsernameRequired          = errors.New("username is required")
	ErrInvalidEmail              = errors.New("invalid email")
	ErrInvalidPhone              = errors.New("invalid phone number - accepts only numbers with at least 10 digits")
	ErrInvalidUsername           = errors.New("invalid username - accepts only letters, numbers, underscores and hyphens, and must be between 3 and 20 characters")
	ErrPasswordRequired          = errors.New("password is required")
	ErrPasswordTooShort          = errors.New("password must be at least 8 characters long")
	ErrInvalidPassword           = errors.New("invalid password - the password must to have special characters, uppercase, lowercase, numbers and be between 8 and 50 characters")
)
