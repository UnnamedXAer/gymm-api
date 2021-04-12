package usecases_test

import (
	"testing"

	"github.com/unnamedxaer/gymm-api/usecases"
)

var (
	userUC usecases.IUserUseCases
	ui     usecases.UserInput = usecases.UserInput{
		Username:     "John Silver",
		EmailAddress: "johnsilver@email.com",
		Password:     "TheSecretestPasswordEver123$%^",
	}
	userID string = "6072d3206144644984a54fa1"
)

func TestCreateUser(t *testing.T) {
	u, _ := userUC.CreateUser(&ui)

	if u.EmailAddress != ui.EmailAddress {
		t.Fatalf("Expected UserRepo.CreateUser to be called and 'EmailAddress' to be '%s', got %s", ui.EmailAddress, u.EmailAddress)
	}
}

func TestGetUserByID(t *testing.T) {
	u, _ := userUC.GetUserByID(userID)

	if u.ID != userID {
		t.Fatalf("Expected UserRepo.GetUserByID to be called and 'ID' to be '%s', got %s", userID, u.ID)
	}
}
