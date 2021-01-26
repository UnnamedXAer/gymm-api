package usecases

import "github.com/unnamedxaer/gymm-api/entities"

type UserGetter interface {
	GetUserByID(id string) (entities.User, error)
}

func GetUserByID(id string) (entities.User, error) {
	return
}
