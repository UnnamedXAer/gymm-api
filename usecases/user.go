package usecases

import "github.com/unnamedxaer/gymm-api/entities"

type UserGetter interface {
	// New creates new error of type EmailAddressInUse
	NewEmailAddressInUse() error
	// GetUserByID it is signature of repo method
	// it's here to not couple both packages
	GetUserByID(id string) (entities.User, error)
	CreateUser(fName, lName, email string, password []byte) (entities.User, error)
}

func GetUserByIDUseCase(userGetter UserGetter) func(id string) (entities.User, error) {
	return func(id string) (entities.User, error) {
		return userGetter.GetUserByID(id)
	}
}

func CreateUserUseCase(userGetter UserGetter) func(fName, lName, email string, password []byte) (entities.User, error) {
	return func(fName, lName, email string, password []byte) (entities.User, error) {
		return userGetter.CreateUser(fName, lName, email, password)
	}
}
