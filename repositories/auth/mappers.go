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

func mapTokenToEntity(data *tokenData) *entities.UserToken {
	return &entities.UserToken{
		ID:        data.ID.Hex(),
		UserID:    data.UserID.Hex(),
		Token:     data.Token,
		CreatedAt: data.CreatedAt.UTC(),
		ExpiresAt: data.ExpiresAt.UTC(),
		Device:    data.Device,
	}
}

func mapTokensToEntities(data []tokenData) []entities.UserToken {

	tokens := make([]entities.UserToken, len(data))

	for i := 0; i < len(data); i++ {
		tokens[i] = *mapTokenToEntity(&data[i])
	}
	return tokens
}
