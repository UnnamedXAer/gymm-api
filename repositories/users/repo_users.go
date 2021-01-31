package users

import (
	"context"
	"fmt"
	"time"

	"github.com/unnamedxaer/gymm-api/entities"
	"github.com/unnamedxaer/gymm-api/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// userData is used only to push data to db
type userData struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username     string             `json:"username,omitempty" bson:"username,omitempty"`
	EmailAddress string             `json:"emailAddress,omitempty" bson:"email_address,omitempty"`
	Password     []byte             `json:"password,omitempty" bson:"password,omitempty"`
	CreatedAt    time.Time          `json:"createdAt,omitempty" bson:"created_at,omitempty"`
}

// GetUserByID retrieves user info from storage
func (r *UserRepository) GetUserByID(id string) (entities.User, error) {
	var ud userData
	var u entities.User
	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		r.l.Error().Msg(err.Error())
		return u, fmt.Errorf("invalid user id: %s", id)
	}

	err = r.col.FindOne(context.Background(), bson.M{"_id": oID}).Decode(&ud)
	// @todo: handle mongo document nil error
	if err != nil {
		r.l.Error().Msg(err.Error())
		return u, err
	}

	u = entities.User{
		ID:           ud.ID.Hex(),
		EmailAddress: ud.EmailAddress,
		Username:     ud.Username,
		CreatedAt:    ud.CreatedAt,
	}
	return u, nil
}

// CreateUser inserts newly registered user into storage
func (r *UserRepository) CreateUser(
	username,
	emailAddress string,
	passwordHash []byte) (u entities.User, err error) {

	now := time.Now().UTC()

	ud := userData{
		Username:     username,
		EmailAddress: emailAddress,
		Password:     passwordHash,
		CreatedAt:    now,
	}

	results, err := r.col.InsertOne(nil, ud)

	if err != nil {
		if repositories.IsDuplicatedError(err) {
			return u, fmt.Errorf("email address already in use")
		}

		return u, err
	}

	u = entities.User{
		results.InsertedID.(primitive.ObjectID).Hex(),
		username,
		emailAddress,
		now,
	}
	return u, nil
}
