package usecases

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/unnamedxaer/gymm-api/entities"
	"golang.org/x/crypto/bcrypt"
)

// UserInput represents data received from req
type UserInput struct {
	// ID           primitive.ObjectID `json:"id" validate:"-"`
	Username     string    `json:"userName" validate:"required,min=2,max=50,printascii"`
	EmailAddress string    `json:"emailAddress" validate:"required,email"`
	Password     string    `json:"password" validate:"required,min=6,max=50,pwd"`
	CreatedAt    time.Time `json:"createdAt"`
}

type UserRepo interface {
	// New creates new error of type EmailAddressInUse
	// NewEmailAddressInUse() error

	// GetUserByID it is signature of repo method
	// it's here to not couple both packages
	GetUserByID(ctx context.Context, id string) (*entities.User, error)
	CreateUser(
		ctx context.Context,
		username,
		email string,
		passwordHash []byte) (*entities.User, error)
}

type UserUseCases struct {
	repo UserRepo
}

type IUserUseCases interface {
	GetUserByID(ctx context.Context, id string) (*entities.User, error)
	CreateUser(ctx context.Context, u *UserInput) (*entities.User, error)
}

func (uc *UserUseCases) GetUserByID(
	ctx context.Context,
	id string) (*entities.User, error) {
	return uc.repo.GetUserByID(ctx, id)
}

func (uc *UserUseCases) CreateUser(
	ctx context.Context,
	u *UserInput) (*entities.User, error) {
	passwordHash, err := hashPassword(u.Password)
	if err != nil {
		return nil, errors.WithMessage(err, "incorrect password, cannot hash")
	}

	return uc.repo.CreateUser(ctx, u.Username, u.EmailAddress, passwordHash)
}

func NewUserUseCases(userRepo UserRepo) IUserUseCases {
	return &UserUseCases{
		repo: userRepo,
	}
}

func hashPassword(pwd string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
}
