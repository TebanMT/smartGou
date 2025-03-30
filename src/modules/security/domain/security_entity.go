package domain

import (
	"regexp"
)

type SecurityEntity struct {
	DailingCode string
	PhoneNumber string
	Email       string
}

func (s *SecurityEntity) FormatPhoneE164() string {
	return s.DailingCode + s.PhoneNumber
}

func (s *SecurityEntity) ValidatePhoneNumber() error {
	if s.PhoneNumber == "" || s.DailingCode == "" {
		return ErrInvalidPhoneOrDailingCode
	}

	if len(s.PhoneNumber) < 10 || len(s.PhoneNumber) > 15 {
		return ErrInvalidPhoneLength
	}

	if len(s.DailingCode) < 1 || len(s.DailingCode) > 5 {
		return ErrInvalidDailingCodeLength
	}

	dailingCodeRegex := regexp.MustCompile(`^\+[1-9]\d{0,2}$`)
	if !dailingCodeRegex.MatchString(s.DailingCode) {
		return ErrInvalidDailingCodeFormat
	}

	phoneNumberRegex := regexp.MustCompile(`^[1-9]\d{9,14}$`)
	if !phoneNumberRegex.MatchString(s.PhoneNumber) {
		return ErrInvalidPhoneFormat
	}

	return nil
}
