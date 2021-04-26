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
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id,omitempty"`
	Token     string             `bson:"token,omitempty"`
	Device    string             `bson:"device,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
	ExpiresAt time.Time          `bson:"expires_at,omitempty"`
}

type refreshTokenData struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id,omitempty"`
	Token     string             `bson:"token,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
	ExpiresAt time.Time          `bson:"expires_at,omitempty"`
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

	uOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.WithMessage(err, "authRepo.GetUserJWTs")
	}

	filter := bson.M{
		"user_id": uOID,
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

	uOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.WithMessage(err, "authRepo.SaveJWT")
	}

	data := tokenData{
		UserID:    uOID,
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

		uOID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			return errors.WithMessage(err, "authRepo.DeleteJWT")
		}
		filter.UserID = uOID
		filter.Device = device
	}

	result, err := repo.tokensCol.DeleteMany(context.Background(), &filter)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil
		}
		return errors.WithMessage(err, "authRepo.DeleteJWT")
	}

	repo.l.Debug().Msgf("authRepo.DeleteJWT userID: %q, deleteCnt: %d", userID, result.DeletedCount)
	return nil
}

func (repo *AuthRepository) SaveRefreshToken(
	userID string,
	token string,
	expiresAt time.Time) (*entities.RefreshToken, error) {
	uOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.WithMessage(err, "authRepo.SaveRefreshToken")
	}

	filter := bson.M{"user_id": uOID}
	update := bson.M{"$set": refreshTokenData{
		UserID:    uOID,
		Token:     token,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	result := repo.refTokensCol.FindOneAndUpdate(context.Background(), &filter, &update, opts)
	if err = result.Err(); err != nil {
		return nil, errors.WithMessage(err, "authRepo.SaveRefreshToken")
	}

	data := &refreshTokenData{}
	err = result.Decode(data)
	if err != nil {
		return nil, errors.WithMessage(err, "authRepo.SaveRefreshToken: decode")
	}

	rt := mapRefreshTokenToEntity(data)
	return rt, nil
}

func (repo *AuthRepository) GetRefreshToken(userID string) (*entities.RefreshToken, error) {
	uOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.WithMessage(err, "authRepo.GetRefreshToken")
	}

	filter := refreshTokenData{
		UserID: uOID,
	}

	result := repo.refTokensCol.FindOne(context.Background(), &filter)
	if err = result.Err(); err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, nil
		}
		return nil, errors.WithMessage(err, "authRepo.GetRefreshToken")
	}

	data := &refreshTokenData{}
	err = result.Decode(data)
	if err != nil {
		return nil, errors.WithMessage(err, "authRepo.GetRefreshToken: decode")
	}

	rt := mapRefreshTokenToEntity(data)
	return rt, nil
}

func (repo *AuthRepository) DeleteRefreshToken(userID string) error {
	uOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.WithMessage(err, "authRepo.DeleteRefreshToken")
	}
	filter := refreshTokenData{
		UserID: uOID,
	}

	result, err := repo.tokensCol.DeleteMany(context.Background(), &filter)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil
		}
		return errors.WithMessage(err, "authRepo.DeleteRefreshToken")
	}

	repo.l.Debug().Msgf("authRepo.DeleteRefreshToken userID: %q, deleteCnt: %d", userID, result.DeletedCount)

	return nil
}
