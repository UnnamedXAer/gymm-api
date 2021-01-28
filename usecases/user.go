package usecases

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/unnamedxaer/gymm-api/entities"
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

func CreateUserUseCase(userGetter UserGetter) func(u *UserData) (entities.User, error) {
	return func(u *UserData) (entities.User, error) {
		// validate
		return userGetter.CreateUser(u.FirstName, u.LastName, u.EmailAddress, u.Password)
	}
}
