package mocks

import (
	"context"
	"strings"
	"time"

	"github.com/unnamedxaer/gymm-api/entities"
	"github.com/unnamedxaer/gymm-api/usecases"
	"golang.org/x/crypto/bcrypt"
)

var (
	ExampleUserToken = entities.UserToken{
		ID:        "609bba6ffc79f3b3bdd9431f",
		UserID:    UserID,
		Token:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjYwNzJkMzIwNjE0NDY0NDk4NGE1NGZhMSIsImNyZWF0ZWRBdCI6MTYyMDQ4Njg2MjkyOTA2NTUwMCwiciI6eyJ0b2tlbiI6IjRiYmI1NDZlLTM0MjEtNGI1Ni1iODI2LWQwYjY0OGM3YjUzYiJ9LCJleHAiOjE2MjA0ODcxNjIsImlhdCI6MTYyMDQ4Njg2Mn0.j8dSyoCoSiKOI1duK9eNAl2LFgHjV0fe8eB58IhxlvI",
		Device:    "PostmanRuntime/7.26.10",
		CreatedAt: Now,
		ExpiresAt: Now.AddDate(1, 0, 0),
	}
	ExampleRefreshToken = entities.RefreshToken{
		ID:        "609bba6ffc79f3b3bdd9431e",
		UserID:    UserID,
		Token:     "4bbb546e-3421-4b56-b826-d0b648c7b53b",
		CreatedAt: Now,
		ExpiresAt: Now.AddDate(1, 0, 0),
	}

	ExampleResetPwdReq = entities.ResetPwdReq{
		ID:           "609ba7fedcc11ebd3051a816",
		EmailAddress: ExampleUser.EmailAddress,
		Status:       entities.ResetPwdStatusNoActionYet,
		ExpiresAt:    Now.Add(time.Minute * 15),
		CreatedAt:    Now,
	}
)

type MockAuthRepo struct{}

func (r *MockAuthRepo) GetUserByEmailAddress(
	ctx context.Context,
	emailAddress string) (*entities.AuthUser, error) {
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

func (r *MockAuthRepo) GetUserByID(ctx context.Context, id string) (*entities.AuthUser, error) {

	if strings.Contains(id, "notfound") {
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

func (r *MockAuthRepo) ChangePassword(
	ctx context.Context,
	userID string,
	newPwd []byte) error {
	return nil
}

func (r *MockAuthRepo) AddResetPasswordRequest(
	ctx context.Context,
	emailaddress string,
	expiresAt time.Time) (*entities.ResetPwdReq, error) {

	if len(emailaddress) == 0 || strings.Contains(emailaddress, "notfound") {
		return nil, usecases.NewErrorRecordNotExists("user")
	}

	out := ExampleResetPwdReq
	out.EmailAddress = emailaddress
	out.ExpiresAt = expiresAt
	return &out, nil
}

func (r *MockAuthRepo) UpdatePasswordForResetRequest(
	ctx context.Context,
	reqID string,
	pwdHash []byte) error {

	if strings.Contains(reqID, "notfound") {
		return usecases.NewErrorRecordNotExists("reset password request")
	}

	if len(reqID) == 0 {
		return usecases.NewErrorInvalidID(reqID, "reset password request")
	}

	return nil
}

func (r *MockAuthRepo) GetUserJWTs(
	ctx context.Context,
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

func (r *MockAuthRepo) SaveJWT(
	ctx context.Context,
	userID string,
	device string,
	token string,
	expiresAt time.Time) (*entities.UserToken, error) {
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

func (r *MockAuthRepo) DeleteJWT(
	ctx context.Context,
	token *entities.UserToken) (int64, error) {

	return 1, nil
}

func (r *MockAuthRepo) SaveRefreshToken(
	ctx context.Context,
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
	ctx context.Context,
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

func (r *MockAuthRepo) DeleteRefreshToken(
	ctx context.Context,
	userID string) (n int64, err error) {

	return 1, nil
}

func (r *MockAuthRepo) DeleteRefreshTokenAndAllTokens(
	ctx context.Context,
	userID string) (n int64, err error) {

	return 2, nil
}
