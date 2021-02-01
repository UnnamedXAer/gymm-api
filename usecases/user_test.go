package usecases

import (
	"os"
	"testing"

	"github.com/unnamedxaer/gymm-api/mocks"
)

var (
	uc IUserUseCases
	ui UserInput = UserInput{
		Username:     "John Silver",
		EmailAddress: "johnsilver@email.com",
		Password:     "TheSecretestPasswordEver123$%^",
	}
	userID string = "dadadada"
)

func TestMain(m *testing.M) {

	var ur UserRepo = mocks.MockUserRepo{}

	uc = NewUserUseCases(ur)

	code := m.Run()
	os.Exit(code)
}

func TestCreateUser(t *testing.T) {
	u, _ := uc.CreateUser(&ui)

	if u.EmailAddress != ui.EmailAddress {
		t.Fatalf("Expected UserRepo.CreateUser to be called and 'EmailAddress' to be '%s', got %s", ui.EmailAddress, u.EmailAddress)
	}
}

func TestGetUserByID(t *testing.T) {
	u, _ := uc.GetUserByID(userID)

	if u.ID != userID {
		t.Fatalf("Expected UserRepo.GetUserByID to be called and 'ID' to be '%s', got %s", userID, u.ID)
	}
}
