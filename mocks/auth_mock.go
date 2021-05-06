package mocks

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/unnamedxaer/gymm-api/entities"
	"golang.org/x/crypto/bcrypt"
)

var (
	ExampleUserToken = entities.UserToken{
		ID:        "ID here",
		UserID:    UserID,
		Token:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjYwNzJkMzIwNjE0NDY0NDk4NGE1NGZhMSIsInIiOnsidG9rZW4iOiJkNDVhNmU3Ni0wNTllLTRlOTEtOTgwZi05YjliM2FmOGIyN2MifSwiZXhwIjo0Nzc1MTMzODUyfQ.D2T803_LhVAfEXBMJ1_z--46ak8CmmXe-772QY1-3o0",
		Device:    "PostmanRuntime/7.26.10",
		CreatedAt: time.Now().UTC(),
		ExpiresAt: time.Now().AddDate(1, 0, 0).UTC(),
	}
	ExampleRefreshToken = entities.RefreshToken{
		ID:        "ID here",
		UserID:    UserID,
		Token:     uuid.NewString(),
		CreatedAt: time.Now().UTC(),
		ExpiresAt: time.Now().AddDate(1, 0, 0).UTC(),
	}
)

type MockAuthRepo struct{}

func (r *MockAuthRepo) GetUserByEmailAddress(emailAddress string) (*entities.AuthUser, error) {
	// mock storage get where ID = id

	if strings.Contains(emailAddress, "notfound") {
		return nil, nil
	}

	pwd, err := bcrypt.GenerateFromPassword(Password, bcrypt.MinCost)
	if err != nil {
		return nil, err
	}

	return &entities.AuthUser{
		User:     ExampleUser,
		Password: pwd,
	}, nil
}

func (r *MockAuthRepo) GetUserJWTs(
	userID string,
	expired entities.ExpireType,
) ([]entities.UserToken, error) {
	if strings.Contains(userID, "notfound") {
		return nil, nil
	}

	expirationTime := ExampleUserToken.ExpiresAt

	switch expired {
	case entities.All:
		break
	case entities.NotExpired:
		break
	case entities.Expired:
		expirationTime = time.Now().AddDate(0, 0, -1)
	}

	return []entities.UserToken{
		{
			ID:        ExampleUser.ID,
			UserID:    userID,
			Token:     ExampleUserToken.Token,
			Device:    ExampleUserToken.Device,
			CreatedAt: ExampleUserToken.CreatedAt,
			ExpiresAt: expirationTime,
		},
	}, nil
}

func (r *MockAuthRepo) SaveJWT(userID string, device string, token string, expiresAt time.Time) (*entities.UserToken, error) {
	out := &entities.UserToken{
		ID:        ExampleUserToken.ID,
		UserID:    userID,
		Token:     token,
		Device:    device,
		CreatedAt: ExampleUserToken.CreatedAt,
		ExpiresAt: expiresAt,
	}

	return out, nil
}

func (r *MockAuthRepo) DeleteJWT(token *entities.UserToken) (int64, error) {

	return 1, nil
}

func (r *MockAuthRepo) SaveRefreshToken(
	userID string,
	token string,
	expiresAt time.Time) (*entities.RefreshToken, error) {
	out := &entities.RefreshToken{
		ID:        ExampleRefreshToken.ID,
		UserID:    userID,
		Token:     token,
		CreatedAt: ExampleRefreshToken.CreatedAt,
		ExpiresAt: expiresAt,
	}

	return out, nil
}

func (r *MockAuthRepo) GetRefreshToken(
	userID string) (*entities.RefreshToken, error) {
	out := &entities.RefreshToken{
		ID:        ExampleRefreshToken.ID,
		UserID:    userID,
		Token:     ExampleRefreshToken.Token,
		CreatedAt: ExampleRefreshToken.CreatedAt,
		ExpiresAt: ExampleRefreshToken.ExpiresAt,
	}

	return out, nil
}

func (r *MockAuthRepo) DeleteRefreshToken(userID string) error {

	return nil
}
