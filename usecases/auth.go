package usecases

import (
	"time"

	"github.com/pkg/errors"
	"github.com/unnamedxaer/gymm-api/entities"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepo interface {
	GetUserByEmailAddress(emailAddress string) (*entities.AuthUser, error)
	SaveJWT(userID string, device string, token string, expiresAt time.Time) (*entities.UserToken, error)
	GetUserJWTs(userID string, expired entities.ExpireType) ([]entities.UserToken, error)
	DeleteJWT(ut *entities.UserToken) (int64, error)
	SaveRefreshToken(userID string, token string, expiresAt time.Time) (*entities.RefreshToken, error)
	GetRefreshToken(userID string) (*entities.RefreshToken, error)
	DeleteRefreshToken(userID string) (n int64, err error)
	DeleteRefreshTokenAndAllTokens(userID string) (n int64, err error)
}

type AuthUsecases struct {
	repo AuthRepo
}

type IAuthUsecases interface {
	// Login checks given credentials against registered users
	Login(u *UserInput) (*entities.User, error)
	// SaveJWT saves jwt for given user and device name
	SaveJWT(userID string, device string, token string, expiresAt time.Time) (*entities.UserToken, error)
	// GetUserJWTs returns user jwt tokens
	// if expired is true it returns only expired tokens
	GetUserJWTs(userID string, expired entities.ExpireType) ([]entities.UserToken, error)
	// DeleteJWT removes jwt token, it returns number of deleted results and error if any
	DeleteJWT(ut *entities.UserToken) (int64, error)
	// SaveRefreshToken creates new or override existing entry in storage with given refresh token for the user.
	SaveRefreshToken(userID string, token string, expiresAt time.Time) (*entities.RefreshToken, error)
	// GetRefreshToken reads refresh token for given user.
	GetRefreshToken(userID string) (*entities.RefreshToken, error)
	// DeleteRefreshToken removes refresh token for given user id
	DeleteRefreshToken(userID string) (n int64, err error)
	// DeleteRefreshTokenAndAllTokens removes all jwt tokens and refresh token for given user
	DeleteRefreshTokenAndAllTokens(userID string) (n int64, err error)
}

type IncorrectCredentialsError struct{}

func (err IncorrectCredentialsError) Error() string {
	return "incorrect credentials"
}

func (au *AuthUsecases) Login(u *UserInput) (*entities.User, error) {
	user, err := au.repo.GetUserByEmailAddress(u.EmailAddress)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(u.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, IncorrectCredentialsError{}
		}

		return nil, errors.WithMessage(err, "usecases.Login")
	}

	return &user.User, nil
}

func (au *AuthUsecases) SaveJWT(userID string, device string, token string, expiresAt time.Time) (*entities.UserToken, error) {
	// @todo: drop previous token(s) for this device
	return au.repo.SaveJWT(userID, device, token, expiresAt)
}

func (au *AuthUsecases) DeleteJWT(token *entities.UserToken) (int64, error) {
	return au.repo.DeleteJWT(token)
}

func (au *AuthUsecases) GetUserJWTs(userID string, expired entities.ExpireType) ([]entities.UserToken, error) {
	return au.repo.GetUserJWTs(userID, expired)
}

func (au *AuthUsecases) SaveRefreshToken(userID string, token string, expiresAt time.Time) (*entities.RefreshToken, error) {
	return au.repo.SaveRefreshToken(userID, token, expiresAt)
}

func (au *AuthUsecases) GetRefreshToken(userID string) (*entities.RefreshToken, error) {
	return au.repo.GetRefreshToken(userID)
}

func (au *AuthUsecases) DeleteRefreshToken(userID string) (n int64, err error) {
	return au.repo.DeleteRefreshToken(userID)
}

func (au *AuthUsecases) DeleteRefreshTokenAndAllTokens(userID string) (n int64, err error) {
	return au.repo.DeleteRefreshTokenAndAllTokens(userID)
}

// NewAuthUsecases creates auth usecases
func NewAuthUsecases(repo AuthRepo) IAuthUsecases {
	return &AuthUsecases{
		repo: repo,
	}
}
