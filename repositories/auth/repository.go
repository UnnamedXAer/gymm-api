package auth

import (
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthRepository struct {
	col *mongo.Collection
	l   *zerolog.Logger
}

func NewRepository(logger *zerolog.Logger, col *mongo.Collection) *AuthRepository {
	return &AuthRepository{
		col: col,
		l:   logger,
	}
}
