package usecases_test

import (
	"context"
	"testing"

	"github.com/unnamedxaer/gymm-api/mocks"
	"github.com/unnamedxaer/gymm-api/usecases"
)

var (
	userUC    usecases.IUserUseCases
	userInput usecases.UserInput = usecases.UserInput{
		Username:     mocks.ExampleUser.Username,
		EmailAddress: mocks.ExampleUser.EmailAddress,
		Password:     string(mocks.Password),
	}
)

func TestCreateUser(t *testing.T) {
	ctx := context.TODO()
	got, _ := userUC.CreateUser(ctx, &userInput)

	if got.EmailAddress != userInput.EmailAddress {
		t.Fatalf("want UserRepo.CreateUser to be called and 'EmailAddress' to be '%s', got %s", userInput.EmailAddress, got.EmailAddress)
	}
}

func TestGetUserByID(t *testing.T) {
	ctx := context.TODO()
	got, _ := userUC.GetUserByID(ctx, mocks.UserID)

	if got.ID != mocks.UserID {
		t.Fatalf("want UserRepo.GetUserByID to be called and 'ID' to be '%s', got %s", mocks.UserID, got.ID)
	}
}
