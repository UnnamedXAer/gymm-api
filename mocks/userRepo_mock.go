package mocks

import (
	"context"
	"strings"
	"time"

	"github.com/unnamedxaer/gymm-api/entities"
	"github.com/unnamedxaer/gymm-api/repositories"
)

var (
	ExampleUser = entities.User{
		ID:           UserID,
		Username:     "John Silver",
		EmailAddress: "johnsilver@email.com",
		CreatedAt:    time.Now(),
	}
)

type MockUserRepo struct {
}

func (ur MockUserRepo) GetUserByID(
	ctx context.Context,
	id string) (*entities.User, error) {
	// mock storage get where ID = id

	if strings.Contains(id, "notfound") {
		return nil, nil
	}

	if strings.Contains(id, "INVALIDID") {
		return nil, repositories.NewErrorInvalidID(id, "user")
	}

	return &entities.User{
		ID: id,
	}, nil
}

func (ur MockUserRepo) CreateUser(
	ctx context.Context,
	username, email string, passwordHash []byte) (*entities.User, error) {
	// mock storage insert new user
	u := ExampleUser
	u.Username = username
	u.EmailAddress = email
	return &u, nil
}
