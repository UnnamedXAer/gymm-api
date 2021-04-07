package trainings

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/unnamedxaer/gymm-api/entities"
	"github.com/unnamedxaer/gymm-api/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type trainingData struct {
	ID        primitive.ObjectID     `bson:"_id,omitempty,required"`
	UserID    primitive.ObjectID     `bson:"user_id,omitempty,required"`
	StartTime time.Time              `bson:"start_time,omitempty,required"`
	EndTime   time.Time              `bson:"end_time,omitempty"`
	Exercises []trainingExerciseData `bson:"exercises,omitempty"`
	Comment   string                 `bson:"comment,omitempty"`
}

type trainingExerciseData struct {
	ID         primitive.ObjectID `bson:"_id,omitempty,required"`
	ExerciseID primitive.ObjectID `bson:"exercise_id,omitempty,required"`
	StartTime  time.Time          `bson:"start_time,omitempty,required"`
	EndTime    time.Time          `bson:"end_time,omitempty"`
	Sets       []trainingSetData  `bson:"sets,omitempty"`
	Comment    string             `bson:"comment,omitempty"`
}
type trainingSetData struct {
	ID   primitive.ObjectID `bson:"_id,omitempty,required"`
	Time time.Time          `bson:"time,omitempty,required"`
	Reps int                `bson:"reps,omitempty,required"`
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

func (r *TrainingRepository) EndTraining(trainingID string, endTime time.Time) (t entities.Training, err error) {
	tOID, err := primitive.ObjectIDFromHex(trainingID)
	if err != nil {
		r.l.Error().Msg(err.Error())
		return t, repositories.NewErrorInvalidID(trainingID)
	}

	returnDoc := options.After
	upsert := false
	options := options.FindOneAndUpdateOptions{
		ReturnDocument: &returnDoc,
		Upsert:         &upsert,
	}

	results := r.col.FindOneAndUpdate(context.Background(), trainingData{ID: tOID}, trainingData{EndTime: endTime}, &options)
	err = results.Err()
	if err != nil {
		msg := fmt.Sprintf("cannot update training with ID %s, error: %v", trainingID, err)
		r.l.Error().Msg(msg)
		if repositories.IsDuplicatedError(err) {
			//
		} else if strings.Contains(err.Error(), "") {
			//
		} else if errors.Is(err, mongo.ErrNoDocuments) {
			return t, repositories.NewErrorNotFoundRecord()
		}
		return t, err
	}

	td := trainingData{}
	err = results.Decode(&td)
	if err != nil {
		r.l.Err(err).Send()
		return t, err
	}

	t.ID = td.ID.Hex()
	t.UserID = td.UserID.Hex()
	t.StartTime = td.StartTime
	t.EndTime = td.EndTime
	t.Exercises = make([]entities.TrainingExercise, len(td.Exercises))
	t.Comment = td.Comment

	for i, exData := range td.Exercises {
		exercise := entities.TrainingExercise{
			ID:         exData.ID.Hex(),
			ExerciseID: exData.ExerciseID.Hex(),
			StartTime:  exData.StartTime,
			EndTime:    exData.EndTime,
			Sets:       make([]entities.TrainingSet, len(exData.Sets)),
			Comment:    exData.Comment,
		}
		for j, setData := range exData.Sets {
			exercise.Sets[j] = entities.TrainingSet{
				ID:   setData.ID.Hex(),
				Time: setData.Time,
				Reps: setData.Reps,
			}
		}

		t.Exercises[i] = exercise
	}

	return t, nil
}

func (r *TrainingRepository) GetStartedTrainings(userID string) (t []entities.Training, err error) {
	oUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		r.l.Error().Msg(err.Error())
		return t, fmt.Errorf("invalid user id: %s", userID)
	}

	cursor, err := r.col.Find(context.TODO(), bson.M{"user_id": oUserID, "start_time": nil})
	if err != nil {
		r.l.Error().Msg(err.Error())
		return t, err
	}

	panic("TrainingRepository - GetStartedTrainings - not implemented yet.")
	for cursor.Next(context.Background()) {
		// t = append(t, entities.Training{
		// 	ID: cursor.,
		// })
	}
	return t, nil
}
