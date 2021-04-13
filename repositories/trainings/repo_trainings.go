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

func (r *TrainingRepository) StartTraining(userID string, startTime time.Time) (t *entities.Training, err error) {
	ouID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return t, fmt.Errorf("start training: %v", err)
	}
	td := trainingData{
		UserID:    ouID,
		StartTime: startTime,
	}
	results, err := r.col.InsertOne(context.TODO(), td)
	if err != nil {
		return t, fmt.Errorf("start training: %v", err)
	}

	t = &entities.Training{
		ID:        results.InsertedID.(primitive.ObjectID).Hex(),
		StartTime: startTime,
		UserID:    userID,
	}

	return t, nil
}

func (r *TrainingRepository) EndTraining(trainingID string, endTime time.Time) (t *entities.Training, err error) {
	tOID, err := primitive.ObjectIDFromHex(trainingID)
	if err != nil {
		r.l.Error().Msg(err.Error())
		return t, repositories.NewErrorInvalidID(trainingID)
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	update := bson.M{"$set": bson.M{"end_time": endTime}}

	results := r.col.FindOneAndUpdate(context.Background(), trainingData{ID: tOID}, update, opts)
	err = results.Err()
	if err != nil {
		if repositories.IsDuplicatedError(err) {
			//
		} else if strings.Contains(err.Error(), "") {
			//
		} else if errors.Is(err, mongo.ErrNoDocuments) {
			return t, repositories.NewErrorNotFoundRecord()
		}
		return t, fmt.Errorf("end training: %v", err)
	}

	td := trainingData{}
	err = results.Decode(&td)
	if err != nil {
		r.l.Err(err).Send()
		return t, fmt.Errorf("end training: %v", err)
	}

	t = &entities.Training{
		ID:        td.ID.Hex(),
		UserID:    td.UserID.Hex(),
		StartTime: td.StartTime,
		EndTime:   td.EndTime,
		Exercises: make([]entities.TrainingExercise, len(td.Exercises)),
		Comment:   td.Comment,
	}

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
		return nil, fmt.Errorf("get user trainings: %v", err)
	}

	filter := bson.M{}
	filter["user_id"] = oUserID
	if started {
		filter["end_time"] = nil
	}

	cursor, err := r.col.Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("get user trainings: %v", err)
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.Background()) {
		var training trainingData
		err = cursor.Decode(&training)
		if err != nil {
			return nil, fmt.Errorf("get user trainings: %v", err)
		}

		t = append(t, mapTrainingToEntity(training))
	}

	if err = cursor.Err(); err != nil {
		return nil, fmt.Errorf("get user trainings: %v", err)
	}

	return t, nil
}

func (r TrainingRepository) AddExercise(trID string, exercise *entities.TrainingExercise) (*entities.TrainingExercise, error) {
	tOID, err := primitive.ObjectIDFromHex(trID)
	if err != nil {
		return nil, repositories.NewErrorInvalidID(trID)
	}

	exOID, err := primitive.ObjectIDFromHex(exercise.ExerciseID)
	if err != nil {
		return nil, repositories.NewErrorInvalidID(exercise.ExerciseID)
	}
	newExerciseData := trainingExerciseData{
		ID:         primitive.NewObjectID(),
		ExerciseID: exOID,
		StartTime:  exercise.StartTime,
		Comment:    exercise.Comment,
	}

	update := bson.M{"$push": bson.M{"exercises": newExerciseData}}
	filter := bson.M{"_id": tOID}

	results, err := r.col.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, fmt.Errorf("add exercise: %v", err)
	}

	if results.ModifiedCount == 0 {
		return nil, fmt.Errorf("add exercise: no documents were modified")
	}

	newExercise := entities.TrainingExercise{
		ID:         newExerciseData.ID.Hex(),
		ExerciseID: newExerciseData.ExerciseID.Hex(),
		StartTime:  newExerciseData.StartTime,
		EndTime:    newExerciseData.EndTime,
		Comment:    newExerciseData.Comment,
		Sets:       mapSetsToEntity(newExerciseData.Sets),
	}
	return &newExercise, nil
}

func (r TrainingRepository) AddSet(teID string, set *entities.TrainingSet) (*entities.TrainingSet, error) {
	teOID, err := primitive.ObjectIDFromHex(teID)
	if err != nil {
		return nil, repositories.NewErrorInvalidID(teID)
	}

	newSetData := trainingSetData{
		ID:   primitive.NewObjectID(),
		Time: set.Time,
		Reps: set.Reps,
	}

	filter := bson.M{
		"exercises._id": teOID,
	}
	// @improvement: check if there is a type safe way to insert nested docs
	update := bson.M{"$push": bson.M{"exercises.$.sets": newSetData}}

	results, err := r.col.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}

	if results.ModifiedCount == 0 {
		return nil, fmt.Errorf("add set: no documents were modified")
	}

	newSet := entities.TrainingSet{
		ID:   newSetData.ID.Hex(),
		Time: newSetData.Time,
		Reps: newSetData.Reps,
	}
	return &newSet, nil
}

func (r TrainingRepository) GetTrainingExercises(id string) ([]entities.TrainingExercise, error) {
	tOID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, repositories.NewErrorInvalidID(id)
	}

	filter := bson.M{"_id": tOID}

	c, err := r.col.Find(context.Background(), filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("get training exercises: %v", err)
	}
	var tr []trainingData
	err = c.All(context.Background(), &tr)
	if err != nil {
		return nil, fmt.Errorf("get training exercises: %v", err)
	}

	te := mapExercisesToEntity(tr[0].Exercises)
	return te, nil
}

func (r TrainingRepository) EndExercise(id string, endTime time.Time) (*entities.TrainingExercise, error) {
	teOID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("end exercise: %v", err)
	}

	filter := bson.M{"exercises._id": teOID}
	update := bson.M{"$set": bson.M{"exercises.$.end_time": endTime.UTC()}}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	result := r.col.FindOneAndUpdate(context.Background(), filter, update, opts)
	if err = result.Err(); err != nil {
		return nil, fmt.Errorf("end exercise: %v", err)
	}

	var td trainingData
	err = result.Decode(&td)
	if err != nil {
		return nil, fmt.Errorf("end exercise: %v", err)
	}
	te := mapExerciseToEntity(td.Exercises[0])
	return &te, nil
}
