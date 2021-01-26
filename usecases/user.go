package usecases

import "github.com/unnamedxaer/gymm-api/entities"

type UserGetter interface {
	// GetUserByID it is signature of repo method
	// it's here to not couple both packages
	GetUserByID(id string) (entities.User, error)
	CreateUser()
}

func GetUserByIDUseCase(userGetter UserGetter) func(id string) (entities.User, error) {
	return func(id string) (entities.User, error) {
		return userGetter.GetUserByID(id)
	}
}
