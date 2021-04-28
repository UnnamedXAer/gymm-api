package users

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/unnamedxaer/gymm-api/entities"
	"github.com/unnamedxaer/gymm-api/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserData is used only to push data to db
type UserData struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username     string             `json:"username,omitempty" bson:"username,omitempty"`
	EmailAddress string             `json:"emailAddress,omitempty" bson:"email_address,omitempty"`
	Password     []byte             `json:"password,omitempty" bson:"password,omitempty"`
	CreatedAt    time.Time          `json:"createdAt,omitempty" bson:"created_at,omitempty"`
}

// GetUserByID retrieves user info from storage
func (r *UserRepository) GetUserByID(id string) (*entities.User, error) {
	var ud UserData
	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.WithMessage(repositories.NewErrorInvalidID(id, "user"), "repo.GetUserByID")
	}

	err = r.col.FindOne(context.Background(), bson.M{"_id": oID}).Decode(&ud)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, nil
		}
		return nil, errors.WithMessage(err, "repo.GetUserByID")
	}

	u := entities.User{
		ID:           ud.ID.Hex(),
		EmailAddress: ud.EmailAddress,
		Username:     ud.Username,
		CreatedAt:    ud.CreatedAt,
	}
	return &u, nil
}

// CreateUser inserts newly registered user into storage
func (r *UserRepository) CreateUser(
	username,
	emailAddress string,
	passwordHash []byte) (*entities.User, error) {

	now := time.Now().UTC()

	ud := UserData{
		Username:     username,
		EmailAddress: emailAddress,
		Password:     passwordHash,
		CreatedAt:    now,
	}

	result, err := r.col.InsertOne(context.Background(), ud)
	if err != nil {
		if repositories.IsDuplicatedError(err) {
			return nil, errors.WithMessage(
				repositories.NewErrorEmailAddressInUse(), "repo.GetUserByID")
		}

		return nil, errors.WithMessage(err, "repo.GetUserByID")
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		r.l.Error().Msgf("repo.CreateUser: id type assertion failed, id: %v", result.InsertedID)
	}

	u := entities.User{
		ID:           id.Hex(),
		Username:     username,
		EmailAddress: emailAddress,
		CreatedAt:    now,
	}
	return &u, nil
}
