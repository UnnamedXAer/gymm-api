package auth

import (
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthRepository struct {
	usersCol     *mongo.Collection
	tokensCol    *mongo.Collection
	refTokensCol *mongo.Collection
	l            *zerolog.Logger
}

func NewRepository(
	logger *zerolog.Logger,
	usersCol, tokensCol, refTokensCol *mongo.Collection) *AuthRepository {
	return &AuthRepository{
		usersCol:     usersCol,
		tokensCol:    tokensCol,
		refTokensCol: refTokensCol,
		l:            logger,
	}
}
