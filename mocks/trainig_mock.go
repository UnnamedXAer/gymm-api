package mocks

import (
	"github.com/unnamedxaer/gymm-api/entities"
	"github.com/unnamedxaer/gymm-api/usecases"
)

// InsertMockUser inserts mocked user to repository with use of the repo functionality
func InsertMockUser(ur usecases.UserRepo) (entities.User, error) {
	return ur.CreateUser(
		"John Silver",
		"johnsilver@email.com",
		[]byte("TheSecretestPasswordEver123$%^"),
	)
}

func StartMockTraining(tr usecases.TrainingRepo) (entities.Training, error) {
	return tr.StartTraining()
}
