package exercises

import (
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
)

type ExerciseRepository struct {
	col *mongo.Collection
	l   *zerolog.Logger
}

func NewRepository(logger *zerolog.Logger, collection *mongo.Collection) *ExerciseRepository {
	return &ExerciseRepository{
		collection,
		logger,
	}
}
