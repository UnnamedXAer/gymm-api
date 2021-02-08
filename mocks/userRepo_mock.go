package mocks

import (
	"strings"
	"time"

	"github.com/unnamedxaer/gymm-api/entities"
	"github.com/unnamedxaer/gymm-api/repositories"
)

type MockUserRepo struct {
}

func (ur MockUserRepo) GetUserByID(id string) (entities.User, error) {
	// mock storage get where ID = id

	if strings.Contains(id, "notfound") {
		return entities.User{}, repositories.NewErrorNotFoundRecord()
	}

	if strings.Contains(id, "INVALIDID") {
		return entities.User{}, repositories.NewErrorInvalidID(id)
	}

	return entities.User{
		ID: id,
	}, nil
}

func (ur MockUserRepo) CreateUser(username, email string, passwordHash []byte) (entities.User, error) {
	// mock storage insert new user

	return entities.User{
		Username:     username,
		EmailAddress: email,
		CreatedAt:    time.Now(),
	}, nil
}
