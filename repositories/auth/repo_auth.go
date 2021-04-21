package auth

import (
	"context"

	"github.com/pkg/errors"
	"github.com/unnamedxaer/gymm-api/entities"
	"github.com/unnamedxaer/gymm-api/repositories/users"
)

func (repo *AuthRepository) GetUserByEmailAddress(emailAddress string) (*entities.AuthUser, error) {
	var ud users.UserData
	// filter := bson.M{"email_address": emailAddress}
	filter := users.UserData{
		EmailAddress: emailAddress,
	}
	err := repo.col.FindOne(context.Background(), &filter).Decode(&ud)
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
