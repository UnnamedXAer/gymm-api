package trainings

import (
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
)

type TrainingRepository struct {
	col *mongo.Collection
	l   *zerolog.Logger
}

func NewRepository(logger *zerolog.Logger, collection *mongo.Collection) *TrainingRepository {
	return &TrainingRepository{
		collection,
		logger,
	}
}
