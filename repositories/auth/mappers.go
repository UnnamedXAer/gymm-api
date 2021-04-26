package auth

import "github.com/unnamedxaer/gymm-api/entities"

func mapRefreshTokenToEntity(data *refreshTokenData) *entities.RefreshToken {
	return &entities.RefreshToken{
		ID:        data.ID.Hex(),
		UserID:    data.UserID.Hex(),
		Token:     data.Token,
		CreatedAt: data.CreatedAt.UTC(),
		ExpiresAt: data.ExpiresAt.UTC(),
	}
}
