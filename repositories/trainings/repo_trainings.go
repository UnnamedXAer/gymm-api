package trainings

import (
	"context"
	"fmt"
	"time"

	"github.com/unnamedxaer/gymm-api/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type trainingData struct {
	ID        primitive.ObjectID `bson:"_id,omitempty,required"`
	UserID    primitive.ObjectID `bson:"user_id,omitempty,required"`
	StartTime time.Time          `bson:"start_time,omitempty,required"`
	EndTime   time.Time          `bson:"end_time,omitempty"`
	Comment   string             `bson:"comment,omitempty"`
}

func (r *TrainingRepository) StartTraining(userID string, startTime time.Time) (t entities.Training, err error) {
	ouID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		r.l.Error().Msg(err.Error())
		return t, fmt.Errorf("invalid user id: %s", userID)
	}
	td := trainingData{
		UserID:    ouID,
		StartTime: startTime,
	}
	results, err := r.col.InsertOne(context.TODO(), td)
	if err != nil {
		r.l.Error().Msg(err.Error())
		return t, err
	}
	return entities.Training{
		ID:        results.InsertedID.(primitive.ObjectID).Hex(),
		StartTime: startTime,
		UserID:    userID,
	}, nil
}

func (r *TrainingRepository) GetStartedTraining(userID string) (t entities.Training, err error) {
	oUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		r.l.Error().Msg(err.Error())
		return t, fmt.Errorf("invalid user id: %s", userID)
	}

	results, err := r.col.InsertOne(context.TODO(), bson.M{"user_id": oUserID, "start_time": nil})
	if err != nil {
		r.l.Error().Msg(err.Error())
		return t, err
	}
	return entities.Training{
		ID:     results.InsertedID.(primitive.ObjectID).Hex(),
		UserID: userID,
	}, nil
}