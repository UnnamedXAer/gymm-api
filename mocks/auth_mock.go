package mocks

import (
	"strings"

	"github.com/unnamedxaer/gymm-api/entities"
	"golang.org/x/crypto/bcrypt"
)

type MockAuthRepo struct{}

func (r *MockAuthRepo) GetUserByEmailAddress(emailAddress string) (*entities.AuthUser, error) {
	// mock storage get where ID = id

	if strings.Contains(emailAddress, "notfound") {
		return nil, nil
	}

	pwd, err := bcrypt.GenerateFromPassword(Password, bcrypt.MinCost)
	if err != nil {
		return nil, err
	}

	return &entities.AuthUser{
		User:     ExampleUser,
		Password: pwd,
	}, nil
}
