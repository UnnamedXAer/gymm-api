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

func (r *TrainingRepository) GetUserTrainings(userID string, started bool) (t []entities.Training, err error) {
	oUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		r.l.Error().Msg(err.Error())
		return nil, fmt.Errorf("invalid user id: %s", userID)
	}

	filter := bson.M{}
	filter["user_id"] = oUserID
	if started {
		filter["end_time"] = nil
	}
	//bson.M{"user_id": oUserID, "end_time": nil}
	cursor, err := r.col.Find(context.Background(), filter)
	if err != nil {
		r.l.Error().Msg(err.Error())
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.Background()) {
		var training trainingData
		err = cursor.Decode(&training)
		if err != nil {
			return nil, err
		}

		t = append(t, mapTrainingToEntity(training))
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return t, nil
}

func (r TrainingRepository) AddExercise(trID string, exercise entities.TrainingExercise) (entities.TrainingExercise, error) {
	tOID, err := primitive.ObjectIDFromHex(trID)
	if err != nil {
		return entities.TrainingExercise{}, repositories.NewErrorInvalidID(trID)
	}

	exOID, err := primitive.ObjectIDFromHex(exercise.ExerciseID)
	if err != nil {
		return entities.TrainingExercise{}, repositories.NewErrorInvalidID(exercise.ExerciseID)
	}
	newExercise := trainingExerciseData{
		ID:         primitive.NewObjectID(),
		ExerciseID: exOID,
		StartTime:  exercise.StartTime,
		Comment:    exercise.Comment,
	}

	update := bson.M{"$push": bson.M{"exercises": newExercise}}
	filter := bson.M{"_id": tOID}

	results, err := r.col.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return entities.TrainingExercise{}, err
	}

	if results.ModifiedCount == 0 {
		return entities.TrainingExercise{}, fmt.Errorf("add exercise: no documents were modified")
	}

	return entities.TrainingExercise{
		ID:         newExercise.ID.Hex(),
		ExerciseID: newExercise.ExerciseID.Hex(),
		StartTime:  newExercise.StartTime,
		EndTime:    newExercise.EndTime,
		Comment:    newExercise.Comment,
		Sets:       mapSetsToEntity(newExercise.Sets),
	}, nil
}

func (r TrainingRepository) AddSet(teID string, set entities.TrainingSet) (entities.TrainingSet, error) {
	tOID, err := primitive.ObjectIDFromHex(teID)
	if err != nil {
		return entities.TrainingSet{}, repositories.NewErrorInvalidID(teID)
	}

	newSet := trainingSetData{
		ID:   primitive.NewObjectID(),
		Time: set.Time,
		Reps: set.Reps,
	}

	update := bson.M{"$push": bson.M{"sets": newSet}}
	filter := bson.M{"_id": tOID}
	// filter does not point to specific exercise
	panic("not implemented yet")
	results, err := r.col.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return entities.TrainingSet{}, err
	}

	if results.ModifiedCount == 0 {
		return entities.TrainingSet{}, fmt.Errorf("add set: no documents were modified")
	}

	return entities.TrainingSet{
		ID:   newSet.ID.Hex(),
		Time: newSet.Time,
		Reps: newSet.Reps,
	}, nil
}

func (r TrainingRepository) GetTrainingExercises(id string) ([]entities.TrainingExercise, error) {
	return []entities.TrainingExercise{}, fmt.Errorf("not implemented yet")
}
func (r TrainingRepository) EndExercise(id string, endTime time.Time) (entities.TrainingExercise, error) {
	return entities.TrainingExercise{}, fmt.Errorf("not implemented yet")
}
