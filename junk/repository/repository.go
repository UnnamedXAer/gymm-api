package repository

import "github.com/unnamedxaer/gymm-api/server/models"

type IRepository interface {
	Initialize(uri string) error
	CreateUser(u *models.User) error
	GetUserById(u *models.User) error
}
