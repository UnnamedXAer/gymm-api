package users

import (
	zerolog "github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	col *mongo.Collection
	l   *zerolog.Logger
}

func NewRepository(logger *zerolog.Logger, collection *mongo.Collection) *UserRepository {
	return &UserRepository{
		col: collection,
		l:   logger,
	}
}
