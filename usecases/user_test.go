package usecases_test

import (
	"testing"

	"github.com/unnamedxaer/gymm-api/mocks"
	"github.com/unnamedxaer/gymm-api/usecases"
)

var (
	userUC usecases.IUserUseCases
	ui     usecases.UserInput = usecases.UserInput{
		Username:     mocks.ExampleUser.Username,
		EmailAddress: mocks.ExampleUser.EmailAddress,
		Password:     "TheSecretestPasswordEver123$%^",
	}
)

func TestCreateUser(t *testing.T) {
	u, _ := userUC.CreateUser(&ui)

	if u.EmailAddress != ui.EmailAddress {
		t.Fatalf("Expected UserRepo.CreateUser to be called and 'EmailAddress' to be '%s', got %s", ui.EmailAddress, u.EmailAddress)
	}
}

func TestGetUserByID(t *testing.T) {
	u, _ := userUC.GetUserByID(mocks.UserID)

	if u.ID != mocks.UserID {
		t.Fatalf("Expected UserRepo.GetUserByID to be called and 'ID' to be '%s', got %s", mocks.UserID, u.ID)
	}
}
