package mocks

import (
	"github.com/unnamedxaer/gymm-api/entities"
	"github.com/unnamedxaer/gymm-api/usecases"
)

// InsertMockUser inserts mocked user to repository with use of the repo functionality
func InsertMockUser(ur usecases.UserRepo) (*entities.User, error) {
	return ur.CreateUser(
		ExampleUser.Username,
		ExampleUser.EmailAddress,
		[]byte(Password),
	)
}

func StartMockTraining(tr usecases.TrainingRepo) (*entities.Training, error) {
	return tr.StartTraining(ExampleTraining.UserID, ExampleTraining.StartTime)
}
