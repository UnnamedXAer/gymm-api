package services

import (
	"fmt"
	"log"

	"github.com/unnamedxaer/gymm-api/repository"
	"github.com/unnamedxaer/gymm-api/server/models"
)

var UService UserService = UserService{}

type UserService struct {
	repo repository.IRepository
}

func (us *UserService) SetRepo(r repository.IRepository) {
	us.repo = r
}
func (us UserService) GetUserById(u *models.User) error {
	log.Println("[UserService.GetUserById] " + fmt.Sprintf("%v", u))
	log.Printf("\n repo: %T, %v", us.repo, us.repo)
	return us.repo.GetUserById(u)
}

func (us UserService) CreateUser(u *models.User) error {
	log.Println("[UserService.CreateUser] " + fmt.Sprintf("%v", u))
	return us.repo.GetUserById(u)
}
