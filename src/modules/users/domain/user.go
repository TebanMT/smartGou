package domain

import (
	"regexp"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                    int
	UserID                uuid.UUID
	Username              *string
	Name                  *string
	LastName              *string
	SecondLastName        *string
	Email                 *string
	DailingCode           *string
	Phone                 *string
	Password              *string
	IsOnboardingCompleted bool
	VerifiedPhone         bool
	VerifiedEmail         bool
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

func (u *User) HashPassword() {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(*u.Password), bcrypt.DefaultCost)
	*u.Password = string(hashedPassword)
}

func (u *User) GenerateUserID() {
	u.UserID = uuid.New()
}

func (u *User) FormatPhoneE164() (string, error) {
	if u.DailingCode == nil || u.Phone == nil {
		return "", ErrInvalidPhone
	}
	formattedPhone := *u.DailingCode + *u.Phone
	re := regexp.MustCompile(`^\+?[1-9]\d{9,14}$`)
	if !re.MatchString(formattedPhone) {
		return "", ErrInvalidPhone
	}
	return formattedPhone, nil
}

type UserValidator interface {
	Validate(user *User) error
}

type EmailValidator struct {
	Next UserValidator
}

func (v *EmailValidator) Validate(user *User) error {
	if user.Email == nil {
		return ErrEmailRequired
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(*user.Email) {
		return ErrInvalidEmail
	}
	if v.Next != nil {
		return v.Next.Validate(user)
	}
	return nil
}

type PhoneValidator struct {
	Next UserValidator
}

func (v *PhoneValidator) Validate(user *User) error {
	if user.Phone == nil {
		return ErrPhoneRequired
	}

	phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{9,14}$`)
	if !phoneRegex.MatchString(*user.Phone) {
		return ErrInvalidPhone
	}

	if v.Next != nil {
		return v.Next.Validate(user)
	}
	return nil
}

type UserNameValidator struct {
	Next UserValidator
}

func (v *UserNameValidator) Validate(user *User) error {
	if user.Username == nil {
		return ErrUsernameRequired
	}

	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]{3,20}$`)
	if !usernameRegex.MatchString(*user.Username) {
		return ErrInvalidUsername
	}

	if v.Next != nil {
		return v.Next.Validate(user)
	}
	return nil
}

type PasswordValidator struct {
	Next UserValidator
}

func (v *PasswordValidator) Validate(user *User) error {
	if user.Password == nil {
		return ErrPasswordRequired
	}

	if len(*user.Password) < 8 {
		return ErrPasswordTooShort
	}

	// Simplified regex that checks for at least one letter, one number, and one special character
	passwordRegex := regexp.MustCompile(`^[a-zA-Z0-9!#_\[\]()$%&? "*]{8,}$`)
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(*user.Password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(*user.Password)
	hasSpecial := regexp.MustCompile(`[!#_\[\]()$%&? "*]`).MatchString(*user.Password)

	if !passwordRegex.MatchString(*user.Password) || !hasLetter || !hasNumber || !hasSpecial {
		return ErrInvalidPassword
	}

	if v.Next != nil {
		return v.Next.Validate(user)
	}
	return nil
}

func BuildUserValidatorChain() UserValidator {
	return &EmailValidator{
		Next: &PhoneValidator{
			Next: &UserNameValidator{
				Next: &PasswordValidator{},
			},
		},
	}
}

func BuildUserValidatorNoPhoneChain() UserValidator {
	return &EmailValidator{
		Next: &UserNameValidator{
			Next: &PasswordValidator{},
		},
	}
}

func BuildUserValidatorEmailAndPasswordChain() UserValidator {
	return &EmailValidator{
		Next: &PasswordValidator{},
	}
}
