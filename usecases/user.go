package usecases

import (
	"time"

	"github.com/unnamedxaer/gymm-api/entities"
	"github.com/unnamedxaer/gymm-api/repositories/users"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserInput represents data received from req
type UserInput struct {
	ID           primitive.ObjectID `json:"id" validate:"-"`
	Username     string             `json:"username" validate:"required,min=2,max=50,printascii"`
	EmailAddress string             `json:"emailAddress" validate:"required,email"`
	Password     string             `json:"password" validate:"required,min=6,max=24,pwd"`
	CreatedAt    time.Time          `json:"createdAt"`
}

type UserRepo interface {
	// New creates new error of type EmailAddressInUse
	// NewEmailAddressInUse() error
	// GetUserByID it is signature of repo method
	// it's here to not couple both packages
	GetUserByID(id string) (entities.User, error)
	CreateUser(username, email string, password string) (entities.User, error)
}

type UserUseCases struct {
	repo UserRepo
}

type IUserUseCases interface {
	GetUserByID(id string) (entities.User, error)
	CreateUser(u *UserInput) (entities.User, error)
}

func (uc *UserUseCases) GetUserByID(id string) (entities.User, error) {
	return uc.repo.GetUserByID(id)
}

func (uc *UserUseCases) CreateUser(u *UserInput) (entities.User, error) {
	return uc.repo.CreateUser(u.Username, u.EmailAddress, u.Password)
}

func NewUserUseCases(userRepo *users.UserRepository) IUserUseCases {
	return &UserUseCases{
		repo: userRepo,
	}
}
