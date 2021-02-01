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

	if strings.Contains(id, "not_found") {
		return entities.User{}, repositories.NewErrorNotFoundRecord()
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
