package mocks

import (
	"fmt"
	"os"

	"github.com/unnamedxaer/gymm-api/entities"
	"github.com/unnamedxaer/gymm-api/usecases"
)

// InsertMockUser inserts mocked user to repository with use of the repo functionality
func InsertMockUser(ur usecases.UserRepo) (entities.User, error) {
	if os.Getenv("ENV") == "test" {
		panic(fmt.Errorf("wrong env, NOT wanted 'test', got '%s'", os.Getenv("ENV")))
	}

	return ur.CreateUser(
		"John Silver",
		"johnsilver@email.com",
		[]byte("TheSecretestPasswordEver123$%^"),
	)
}

func StartMockTraining(tr usecases.TrainingRepo) (entities.Training, error) {
	return tr.StartTraining()
}
