package users

import (
	"context"
	"fmt"
	"time"

	"github.com/unnamedxaer/gymm-api/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userData struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName    string             `json:"firstName,omitempty" bson:"firstName,omitempty"`
	LastName     string             `json:"lastName,omitempty" bson:"lastName,omitempty"`
	EmailAddress string             `json:"emailAddress,omitempty" bson:"emailAddress,omitempty"`
	Password     string             `json:"password,omitempty" bson:"password,omitempty"`
	CreatedAt    time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
}

// Initialize(uri string) error
// CreateUser(u *models.User) error
// GetUserById(u *models.User) error

func (r *UserRepository) GetUserById(id string) (entities.User, error) {
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
