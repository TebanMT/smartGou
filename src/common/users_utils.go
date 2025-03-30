package common

import (
	"strings"

	commonDomain "github.com/TebanMT/smartGou/src/common/domain"
	userDomain "github.com/TebanMT/smartGou/src/modules/users/domain"
)

type CheckType struct {
	ExistsFunc func(commonDomain.Transaction, string) (*userDomain.User, error)
	Value      *string
	Err        error
}

func ValidateUser(user *userDomain.User) error {
	userValidator := userDomain.BuildUserValidatorChain()
	return userValidator.Validate(user)
}

func ValidateUserNoPhone(user *userDomain.User) error {
	userValidator := userDomain.BuildUserValidatorNoPhoneChain()
	return userValidator.Validate(user)
}

func ValidateUserEmailAndPassword(user *userDomain.User) error {
	userValidator := userDomain.BuildUserValidatorEmailAndPasswordChain()
	return userValidator.Validate(user)
}

func checkUserExistence(tx commonDomain.Transaction, checks []CheckType) (*userDomain.User, error) {
	for _, check := range checks {
		if check.Value == nil {
			continue
		}
		user, err := check.ExistsFunc(tx, *check.Value)
		if err != nil {
			return nil, err
		}
		if user != nil {
			return user, check.Err
		}
	}
	return nil, nil
}

func CheckUserExistenceByPhone(tx commonDomain.Transaction, userRepository userDomain.UserRepository, user *userDomain.User) (*userDomain.User, error) {
	checks := []CheckType{
		{userRepository.ExistsUserByPhone, user.Phone, userDomain.ErrPhoneAlreadyExists},
	}
	return checkUserExistence(tx, checks)
}

func CheckUserExistenceByEmail(tx commonDomain.Transaction, userRepository userDomain.UserRepository, user *userDomain.User) (*userDomain.User, error) {
	email := strings.ToLower(*user.Email)
	checks := []CheckType{
		{userRepository.ExistsUserByEmail, &email, userDomain.ErrEmailAlreadyExists},
	}
	return checkUserExistence(tx, checks)
}

func CheckUserExistenceByPhoneAndEmail(tx commonDomain.Transaction, userRepository userDomain.UserRepository, user *userDomain.User) (*userDomain.User, error) {
	email := strings.ToLower(*user.Email)
	checks := []CheckType{
		{userRepository.ExistsUserByPhone, user.Phone, userDomain.ErrPhoneAlreadyExists},
		{userRepository.ExistsUserByEmail, &email, userDomain.ErrEmailAlreadyExists},
	}
	return checkUserExistence(tx, checks)
}
