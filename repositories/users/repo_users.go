package users

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/unnamedxaer/gymm-api/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type userData struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName    string             `json:"firstName,omitempty" bson:"first_name,omitempty"`
	LastName     string             `json:"lastName,omitempty" bson:"last_name,omitempty"`
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
	if err != nil {
		r.l.Error().Msg(err.Error())
		return u, err
	}

	u = entities.User{
		ID:           ud.ID.Hex(),
		EmailAddress: ud.EmailAddress,
		FirstName:    ud.FirstName,
		LastName:     ud.LastName,
		CreatedAt:    ud.CreatedAt,
	}
	return u, nil
}

// CreateUser inserts newly registered user into storage
func (r *UserRepository) CreateUser(firstName, lastName, emailAddress string, passwordHash []byte) (u entities.User, err error) {
	password, err := hashPassword("TheSecretestPasswordEver")
	if err != nil {
		return u, errors.New("incorrect password, cannot hash")
	}
	now := time.Now().UTC()

	ud := userData{
		FirstName:    firstName,
		LastName:     lastName,
		EmailAddress: emailAddress,
		Password:     password,
		CreatedAt:    now,
	}

	panic("todo: ensure email unique")
	results, err := r.col.InsertOne(nil, ud)
	if err != nil {
		return u, err
	}

	return entities.User{
		results.InsertedID.(primitive.ObjectID).Hex(),
		firstName,
		lastName,
		emailAddress,
		now,
	}, nil
}
func hashPassword(pwd string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
}
