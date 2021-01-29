package usecases

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/unnamedxaer/gymm-api/entities"
	"github.com/unnamedxaer/gymm-api/repositories/users"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserData represents data received from req
type UserData struct {
	ID           primitive.ObjectID `json:"id" validate:""`
	FirstName    string             `json:"firstName" validate:"required"`
	LastName     string             `json:"lastName" validate:"required"`
	EmailAddress string             `json:"emailAddress" validate:"required,email"`
	Password     []byte             `json:"password" validate:"required,password"`
	CreatedAt    time.Time          `json:"createdAt" validate:"required,time"`
}

var validate *validator.Validate = validator.New()

type UserRepo interface {
	// New creates new error of type EmailAddressInUse
	// NewEmailAddressInUse() error
	// GetUserByID it is signature of repo method
	// it's here to not couple both packages
	GetUserByID(id string) (entities.User, error)
	CreateUser(fName, lName, email string, password []byte) (entities.User, error)
}

type UserUseCases struct {
	repo UserRepo
}

type IUserUseCases interface {
	GetUserByID(id string) (entities.User, error)
	CreateUser(u *UserData) (entities.User, error)
}

func (uc *UserUseCases) GetUserByID(id string) (entities.User, error) {
	return uc.repo.GetUserByID(id)
}

func (uc *UserUseCases) CreateUser(u *UserData) (entities.User, error) {
	// validate
	return uc.repo.CreateUser(u.FirstName, u.LastName, u.EmailAddress, u.Password)
}

func NewUserUseCases(userRepo *users.UserRepository) IUserUseCases {
	return &UserUseCases{
		repo: userRepo,
	}
}
