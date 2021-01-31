package usecases

import (
	"os"
	"testing"
	"time"

	"github.com/unnamedxaer/gymm-api/entities"
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

type MockUserRepo struct {
}

func (ur MockUserRepo) GetUserByID(id string) (entities.User, error) {
	// mock storage get where ID = id
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

func TestMain(m *testing.M) {

	var ur UserRepo = MockUserRepo{}

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
