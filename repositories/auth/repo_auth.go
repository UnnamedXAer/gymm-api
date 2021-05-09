package auth

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/unnamedxaer/gymm-api/entities"
	"github.com/unnamedxaer/gymm-api/repositories/users"
	"github.com/unnamedxaer/gymm-api/usecases"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func (repo *AuthRepository) GetUserByEmailAddress(
	ctx context.Context,
	emailAddress string) (*entities.AuthUser, error) {

	if emailAddress == "" {
		return nil, errors.WithMessage(
			errors.Errorf("empty email address"), "GetUserByEmailAddress")
	}

	var ud users.UserData

	filter := users.UserData{
		EmailAddress: emailAddress,
	}

	err := repo.usersCol.FindOne(ctx, &filter).Decode(&ud)
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

func (repo *AuthRepository) GetUserByID(ctx context.Context, id string) (*entities.AuthUser, error) {
	uOID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.WithMessage(
			usecases.NewErrorInvalidID(id, "user"), "authRepo.GetUserByID")
	}

	filter := users.UserData{
		ID: uOID,
	}

	var ud users.UserData
	err = repo.usersCol.FindOne(ctx, &filter).Decode(&ud)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, nil
		}
		return nil, errors.WithMessage(err, "authRepo.GetUserByID")
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

func (repo *AuthRepository) ChangePassword(ctx context.Context, userID string, newPwd []byte) error {
	uOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.WithMessage(
			usecases.NewErrorInvalidID(userID, "user"), "authRepo.ChangePassword")
	}

	filter := users.UserData{
		ID: uOID,
	}

	update := bson.M{"$set": bson.M{"password": newPwd}}

	result, err := repo.usersCol.UpdateOne(ctx, &filter, &update)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return errors.WithMessage(
				errors.New("no record has been updated"), "authRepo.ChangePassword")
		}
		return errors.WithMessage(err, "authRepo.ChangePassword")
	}

	if result.MatchedCount == 0 {
		return errors.WithMessage(
			errors.New("no record has been updated"), "authRepo.ChangePassword")
	}

	return nil
}

func (repo *AuthRepository) GetUserJWTs(
	ctx context.Context,

	userID string,
	expired entities.ExpireType,
) ([]entities.UserToken, error) {

	uOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.WithMessage(err, "authRepo.GetUserJWTs")
	}

	filter := bson.M{
		"user_id": uOID,
	}

	switch expired {
	case entities.All:
		break
	case entities.NotExpired:
		filter["expires_at"] = bson.M{"$gt": time.Now()}
	case entities.Expired:
		filter["expires_at"] = bson.M{"$lte": time.Now()}
	}

	cursor, err := repo.tokensCol.Find(ctx, &filter)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, nil
		}
		return nil, errors.WithMessage(err, "authRepo.GetUserJWTs")
	}

	data := []tokenData{}

	err = cursor.All(ctx, &data)
	if err != nil {
		return nil, errors.WithMessage(err, "authRepo.GetUserJWTs")
	}

	tokens := mapTokensToEntities(data)

	return tokens, nil
}

func (repo *AuthRepository) SaveJWT(
	ctx context.Context,
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

	result, err := repo.tokensCol.InsertOne(ctx, &data)
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

func (repo *AuthRepository) DeleteJWT(
	ctx context.Context,
	ut *entities.UserToken) (int64, error) {
	var filter tokenData
	if ut.ID != "" {
		tokenOID, err := primitive.ObjectIDFromHex(ut.ID)
		if err != nil {
			return 0, errors.WithMessage(
				usecases.NewErrorInvalidID(ut.ID, "token"), "authRepo.DeleteJWT")
		}
		filter.ID = tokenOID
	}
	if ut.Token != "" {
		filter.Token = ut.Token
	}
	if ut.UserID != "" {
		uOID, err := primitive.ObjectIDFromHex(ut.UserID)
		if err != nil {
			return 0, errors.WithMessage(
				usecases.NewErrorInvalidID(ut.ID, "token - user"),
				"authRepo.DeleteJWT")
		}
		filter.UserID = uOID
	}
	if ut.Device != "" {
		filter.Device = ut.Device
	}

	result, err := repo.tokensCol.DeleteMany(ctx, &filter)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return result.DeletedCount, nil
		}
		return result.DeletedCount, errors.WithMessage(err, "authRepo.DeleteJWT")
	}

	repo.l.Debug().Msgf(
		"authRepo.DeleteJWT filter: %v, deleteCnt: %d", filter, result.DeletedCount)
	return result.DeletedCount, nil
}

func (repo *AuthRepository) SaveRefreshToken(
	ctx context.Context,
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

	result := repo.refTokensCol.FindOneAndUpdate(ctx, &filter, &update, opts)
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

func (repo *AuthRepository) GetRefreshToken(
	ctx context.Context,
	userID string) (*entities.RefreshToken, error) {
	uOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.WithMessage(err, "authRepo.GetRefreshToken")
	}

	filter := refreshTokenData{
		UserID: uOID,
	}

	result := repo.refTokensCol.FindOne(ctx, &filter)
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

func (repo *AuthRepository) DeleteRefreshToken(
	ctx context.Context,
	userID string) (n int64, err error) {
	uOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return 0, errors.WithMessage(err, "authRepo.DeleteRefreshToken")
	}
	filter := refreshTokenData{
		UserID: uOID,
	}

	result, err := repo.tokensCol.DeleteMany(ctx, &filter)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return result.DeletedCount, nil
		}
		return result.DeletedCount, errors.WithMessage(err,
			"authRepo.DeleteRefreshToken")
	}

	repo.l.Debug().Msgf(
		"authRepo.DeleteRefreshToken userID: %q, deleteCnt: %d",
		userID, result.DeletedCount)

	return result.DeletedCount, nil
}

func (repo *AuthRepository) DeleteRefreshTokenAndAllTokens(
	ctx context.Context,
	userID string) (n int64, err error) {

	uOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		if err != nil {
			return 0, errors.WithMessage(usecases.NewErrorInvalidID(userID, "user"),
				"authRepo.DeleteRefreshTokenAndAllTokens")
		}
	}
	filter := bson.M{
		"user_id": uOID,
	}

	cb := func(sessCtx mongo.SessionContext) (interface{}, error) {

		tResult, err := repo.tokensCol.DeleteMany(sessCtx, &filter)
		if err != nil {
			if err.Error() != "mongo: no documents in result" {
				return tResult.DeletedCount, errors.WithMessage(err,
					"authRepo.DeleteRefreshTokenAndAllTokens: token:")
			}
		}

		rtResults, err := repo.refTokensCol.DeleteMany(sessCtx, &filter)
		if err != nil {
			if err.Error() != "mongo: no documents in result" {
				return rtResults.DeletedCount + tResult.DeletedCount, errors.WithMessage(
					err, "authRepo.DeleteRefreshTokenAndAllTokens: refresh token:")
			}
		}

		return rtResults.DeletedCount + tResult.DeletedCount, nil
	}

	session, err := repo.refTokensCol.Database().Client().StartSession()
	if err != nil {
		return 0, errors.WithMessage(err, "authRepo.DeleteRefreshTokenAndAllTokens")
	}
	defer session.EndSession(ctx)

	results, err := session.WithTransaction(ctx, cb)
	if err != nil {
		return 0, err
	}

	n, ok := results.(int64)
	if !ok {
		repo.l.Debug().Msgf(
			"authRepo.DeleteRefreshTokenAndAllTokens: could not assert, results: %v",
			results)
	} else {
		repo.l.Debug().Msgf(
			"authRepo.DeleteRefreshTokenAndAllTokens userID: %q, deleteCnt: %d",
			userID, n)
	}

	return n, nil
}
