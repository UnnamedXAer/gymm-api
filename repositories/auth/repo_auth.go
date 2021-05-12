package auth

import (
	"context"
	"fmt"
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

type resetPwdData struct {
	ID           primitive.ObjectID      `bson:"_id,omitempty"`
	EmailAddress string                  `bson:"email_address,omitempty"`
	Status       entities.ResetPwdStatus `bson:"status,omitempty"`
	ExpiresAt    time.Time               `bson:"expires_at,omitempty"`
	CreatedAt    time.Time               `bson:"created_at,omitempty"`
	UpdatedAt    time.Time               `bson:"updated_at,omitempty"`
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

func (repo *AuthRepository) AddResetPasswordRequest(ctx context.Context, emailaddress string, expiresAt time.Time) (*entities.ResetPwdReq, error) {
	if expiresAt.Before(time.Now()) {
		return nil, errors.New("expiration time from the past")
	}

	cb := func(sessCtx mongo.SessionContext) (interface{}, error) {
		userFilter := bson.M{
			"email_address": emailaddress,
		}
		countRes, err := repo.usersCol.CountDocuments(sessCtx, userFilter)
		if err != nil {
			return nil, errors.WithMessage(
				err, "add reset password request")

		}
		if countRes == 0 {
			return nil, errors.WithMessage(
				usecases.NewErrorRecordNotExists("user"), "add reset password request")
		}

		insert := resetPwdData{
			EmailAddress: emailaddress,
			Status:       entities.ResetPwdStatusNoActionYet,
			ExpiresAt:    expiresAt,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		result, err := repo.resetPwdCol.InsertOne(sessCtx, &insert)
		if err != nil {
			if err.Error() == "mongo: no documents in result" {
				return nil, errors.New("no record has been created")
			}
			return nil, err
		}

		insertedID, ok := result.InsertedID.(primitive.ObjectID)
		if !ok {
			return nil, errors.New("ID assert failed")
		}

		filter := bson.M{"$and": bson.A{
			bson.M{"status": entities.ResetPwdStatusNoActionYet},
			bson.M{"email_address": emailaddress},
			bson.M{"_id": bson.M{"$ne": insertedID}},
		}}

		update := bson.M{"$set": resetPwdData{
			Status:    entities.ResetPwdStatusCanceled,
			UpdatedAt: time.Now(),
		}}

		// @todo: move before insert and remove $ne condition
		updateResult, err := repo.resetPwdCol.UpdateMany(sessCtx, &filter, &update)
		if err != nil {
			if err.Error() != "mongo: no documents in result" {
				return nil, errors.WithMessage(err,
					"could not cancell previous requests")
			}
		}

		if updateResult.ModifiedCount > 0 {
			repo.l.Debug().Msgf(
				"add reset password request: cancelled %d requests for %s",
				updateResult.ModifiedCount,
				emailaddress)
		}

		insert.ID = insertedID
		return &insert, nil
	}

	session, err := repo.resetPwdCol.Database().Client().StartSession()
	if err != nil {
		return nil, errors.WithMessagef(err, "add reset password request: start session")
	}
	defer session.EndSession(ctx)

	transactionResult, err := session.WithTransaction(ctx, cb)
	if err != nil {
		return nil, errors.WithMessagef(err, "add reset password request")
	}

	data, ok := transactionResult.(*resetPwdData)
	if !ok {
		return nil, fmt.Errorf("add reset password request: session results assertion not of type *resetPwdData")
	}

	resetPwdReq := entities.ResetPwdReq{
		ID:           data.ID.Hex(),
		EmailAddress: data.EmailAddress,
		Status:       data.Status,
		ExpiresAt:    data.ExpiresAt,
		CreatedAt:    data.CreatedAt,
	}

	return &resetPwdReq, nil
}

func (repo *AuthRepository) UpdatePasswordForResetRequest(ctx context.Context, reqID string, pwdHash []byte) error {

	if len(pwdHash) == 0 {
		return errors.WithMessage(fmt.Errorf("missing password"), "reset password request")
	}

	reqOID, err := primitive.ObjectIDFromHex(reqID)
	if err != nil {
		return errors.WithMessage(
			usecases.NewErrorInvalidID(reqID, "reset password request"), "reset password request")
	}

	cb := func(sessCtx mongo.SessionContext) (interface{}, error) {
		reqFilter := bson.M{
			"_id": reqOID,
		}
		result := repo.usersCol.FindOne(sessCtx, reqFilter)
		if err = result.Err(); err != nil {
			if err.Error() == "mongo: no documents in result" {
				return nil, usecases.NewErrorRecordNotExists("reset password request")
			}
			return nil, err
		}

		req := resetPwdData{}

		err := result.Decode(&req)
		if err = result.Err(); err != nil {
			return nil, err
		}

		if req.ExpiresAt.Before(time.Now()) {
			return nil, fmt.Errorf("reset password request expired")
		}

		filter := bson.M{"$and": bson.A{
			bson.M{"email_address": req.EmailAddress},
			bson.M{"status": entities.ResetPwdStatusNoActionYet},
			// @i: not necessarily needed as we would override it in next lines
			bson.M{"_id": bson.M{"$ne": req.ID}},
		}}

		update := bson.M{"$set": resetPwdData{
			Status:    entities.ResetPwdStatusCanceled,
			UpdatedAt: time.Now(),
		}}

		updateResult, err := repo.resetPwdCol.UpdateMany(sessCtx, &filter, &update)
		if err != nil {
			if err.Error() != "mongo: no documents in result" {
				return nil, errors.WithMessage(err,
					"could not cancell previous requests")
			}
		}

		if updateResult.ModifiedCount > 0 {
			repo.l.Debug().Msgf(
				"add reset password request: cancelled %d requests for %s",
				updateResult.ModifiedCount,
				req.EmailAddress)
		}

		update = bson.M{"$set": resetPwdData{
			Status:    entities.ResetPwdStatusCompleted,
			UpdatedAt: time.Now(),
		}}

		updateResult, err = repo.resetPwdCol.UpdateOne(sessCtx, &filter, &update)
		if err != nil {
			return nil, errors.WithMessage(err,
				"could not update")
		}

		if updateResult.MatchedCount == 0 {
			return nil, errors.WithMessage(
				usecases.NewErrorRecordNotExists("reset password request"),
				"could not update")
		}

		// repo.usersCol.UpdateOne()
		return nil, fmt.Errorf("not implemented yet")

		return nil, nil
	}

	session, err := repo.resetPwdCol.Database().Client().StartSession()
	if err != nil {
		return errors.WithMessagef(err, "reset password request: start session")
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, cb)
	if err != nil {
		return errors.WithMessagef(err, "reset password request")
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
