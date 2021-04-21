package usecases

import (
	"github.com/pkg/errors"
	"github.com/unnamedxaer/gymm-api/entities"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepo interface {
	GetUserByEmailAddress(emailAddress string) (*entities.AuthUser, error)
}

type AuthUsecases struct {
	repo AuthRepo
}

type IAuthUsecases interface {
	Login(u *UserInput) (*entities.User, error)
}

type IncorrectCredentialsError struct{}

func (err IncorrectCredentialsError) Error() string {
	return "incorrect credentials"
}

// Login checks given credentials against registered users
func (au *AuthUsecases) Login(u *UserInput) (*entities.User, error) {
	user, err := au.repo.GetUserByEmailAddress(u.EmailAddress)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(u.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, IncorrectCredentialsError{}
		}

		return nil, errors.WithMessage(err, "usecases.Login")
	}

	return &user.User, nil
}

// NewAuthUsecases creates auth usecases
func NewAuthUsecases(repo AuthRepo) IAuthUsecases {
	return &AuthUsecases{
		repo: repo,
	}
}
