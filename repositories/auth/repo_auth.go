package auth

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/unnamedxaer/gymm-api/entities"
	"github.com/unnamedxaer/gymm-api/repositories/users"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type tokenData struct {
	ID        string    `bson:"_id,omitempty"`
	UserID    string    `bson:"user_id"`
	Token     string    `bson:"token"`
	Device    string    `bson:"device"`
	CreatedAt time.Time `bson:"created_at"`
	ExpiresAt time.Time `bson:"expires_at"`
}

type refreshTokenData struct {
	ID        string    `bson:"_id,omitempty"`
	UserID    string    `bson:"user_id"`
	Token     string    `bson:"token"`
	CreatedAt time.Time `bson:"created_at"`
	ExpiresAt time.Time `bson:"expires_at"`
}

func (repo *AuthRepository) GetUserByEmailAddress(emailAddress string) (*entities.AuthUser, error) {
	var ud users.UserData
	// filter := bson.M{"email_address": emailAddress}
	filter := users.UserData{
		EmailAddress: emailAddress,
	}
	err := repo.usersCol.FindOne(context.Background(), &filter).Decode(&ud)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, nil
		}
		return nil, errors.WithMessage(err, "authRepo.GetUserByEmailAddress")
	}

	u := &entities.AuthUser{
		User: entities.User{
			ID:           ud.ID.Hex(),
			EmailAddress: ud.EmailAddress,
			Username:     ud.Username,
			CreatedAt:    ud.CreatedAt,
		},
		Password: ud.Password,
	}
	return u, nil
}

func (repo *AuthRepository) GetUserJWTs(
	userID string,
	expired bool,
) ([]entities.UserToken, error) {

	filter := bson.M{
		"user_id": userID,
	}

	if expired {
		filter["expires_at"] = bson.M{"$lte": time.Now()}
	}

	cursor, err := repo.tokensCol.Find(context.Background(), &filter)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, nil
		}
		return nil, errors.WithMessage(err, "authRepo.GetUserJWTs")
	}

	var tokens []entities.UserToken

	err = cursor.All(context.Background(), &tokens)
	if err != nil {
		return nil, errors.WithMessage(err, "authRepo.GetUserJWTs")
	}

	return tokens, nil
}

func (repo *AuthRepository) SaveJWT(
	userID string,
	device string,
	token string,
	expiresAt time.Time) (*entities.UserToken, error) {

	data := tokenData{
		UserID:    userID,
		Token:     token,
		Device:    device,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}

	result, err := repo.tokensCol.InsertOne(context.Background(), &data)
	if err != nil {
		return nil, errors.WithMessage(err, "authRepo.SaveJWT")
	}

	out := &entities.UserToken{
		ID:        result.InsertedID.(primitive.ObjectID).Hex(),
		UserID:    userID,
		Token:     token,
		Device:    device,
		CreatedAt: data.CreatedAt,
		ExpiresAt: data.ExpiresAt,
	}

	return out, nil
}

func (repo *AuthRepository) DeleteJWT(userID string, device string, token string) error {

	var filter tokenData
	if token != "" {
		filter.Token = token
	} else {
		filter.UserID = userID
		filter.Device = device
	}

	_, err := repo.tokensCol.DeleteMany(context.Background(), &filter)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil
		}
		return errors.WithMessage(err, "authRepo.DeleteJWT")
	}
	return nil
}

func (repo *AuthRepository) SaveRefreshToken(
	userID string,
	token string,
	expiresAt time.Time) (*entities.RefreshToken, error) {
	var err error
	data := refreshTokenData{
		UserID:    userID,
		Token:     token,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}

	opts := options.FindOneAndUpdate().SetUpsert(true)

	result := repo.refTokensCol.FindOneAndUpdate(context.Background(), &data, opts)
	if err = result.Err(); err != nil {
		return nil, errors.WithMessage(err, "authRepo.SaveRefreshToken")
	}

	var rt *entities.RefreshToken

	err = result.Decode(rt)
	if err != nil {
		return nil, errors.WithMessage(err, "authRepo.SaveRefreshToken: decode")
	}

	return rt, nil
}

func (repo *AuthRepository) GetRefreshToken(userID string) (*entities.RefreshToken, error) {
	var err error
	filter := refreshTokenData{
		UserID: userID,
	}

	result := repo.refTokensCol.FindOne(context.Background(), &filter)
	if err = result.Err(); err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, nil
		}
		return nil, errors.WithMessage(err, "authRepo.GetRefreshToken")
	}

	var rt *entities.RefreshToken
	err = result.Decode(rt)
	if err != nil {
		return nil, errors.WithMessage(err, "authRepo.GetRefreshToken: decode")
	}

	return rt, nil
}

func (repo *AuthRepository) DeleteRefreshToken(userID string) error {

	filter := refreshTokenData{
		UserID: userID,
	}

	result, err := repo.tokensCol.DeleteMany(context.Background(), &filter)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil
		}
		return errors.WithMessage(err, "authRepo.DeleteRefreshToken")
	}

	if result.DeletedCount > 1 {
		repo.l.Error().Msgf("authRepo.DeleteRefreshToken: user: %q had %d refresh tokens", userID, result.DeletedCount)
	}

	return nil
}
