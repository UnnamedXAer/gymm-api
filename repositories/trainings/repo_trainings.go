package trainings

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
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
	CreatedAt time.Time              `bson:"created_at,omitempty,required"`
}

type trainingExerciseData struct {
	ID         primitive.ObjectID `bson:"_id,omitempty,required"`
	ExerciseID primitive.ObjectID `bson:"exercise_id,omitempty,required"`
	StartTime  time.Time          `bson:"start_time,omitempty,required"`
	EndTime    time.Time          `bson:"end_time,omitempty"`
	Sets       []trainingSetData  `bson:"sets,omitempty"`
	Comment    string             `bson:"comment,omitempty"`
	CreatedAt  time.Time          `bson:"created_at,omitempty,required"`
}

type trainingSetData struct {
	ID        primitive.ObjectID `bson:"_id,omitempty,required"`
	Time      time.Time          `bson:"time,omitempty,required"`
	Reps      int                `bson:"reps,omitempty,required"`
	CreatedAt time.Time          `bson:"created_at,omitempty,required"`
}

func (r *TrainingRepository) GetTrainingByID(id string) (*entities.Training, error) {
	oTID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.WithMessage(repositories.NewErrorInvalidID(id, "training"), "get training by id")
	}

	filter := bson.M{"_id": oTID}

	result := r.col.FindOne(context.Background(), filter)
	if err = result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, fmt.Errorf("get training by id: %v", err)
	}

	td := trainingData{}
	err = result.Decode(&td)
	if err != nil {
		return nil, fmt.Errorf("get training by id: %v", err)
	}

	t := mapTrainingToEntity(&td)

	return t, nil
}

func (r *TrainingRepository) StartTraining(userID string, startTime time.Time) (*entities.Training, error) {
	ouID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.WithMessage(repositories.NewErrorInvalidID(userID, "user"), "start training")
	}
	td := trainingData{
		UserID:    ouID,
		StartTime: startTime,
		CreatedAt: time.Now(),
	}
	results, err := r.col.InsertOne(context.TODO(), td)
	if err != nil {
		return nil, fmt.Errorf("start training: %v", err)
	}

	t := &entities.Training{
		ID:        results.InsertedID.(primitive.ObjectID).Hex(),
		StartTime: startTime,
		UserID:    userID,
		CreatedAt: td.CreatedAt,
	}

	return t, nil
}

func (r *TrainingRepository) EndTraining(id string, endTime time.Time) (*entities.Training, error) {
	tOID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.WithMessage(repositories.NewErrorInvalidID(id, "training"), "end training")
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
			return nil, nil
		}
		return nil, fmt.Errorf("end training: %v", err)
	}

	td := trainingData{}
	err = results.Decode(&td)
	if err != nil {
		r.l.Err(err).Send()
		return nil, fmt.Errorf("end training: %v", err)
	}

	t := mapTrainingToEntity(&td)

	return t, nil
}

func (r *TrainingRepository) GetUserTrainings(userID string, started bool) ([]entities.Training, error) {
	oUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.WithMessage(repositories.NewErrorInvalidID(userID, "user"), "get user trainings")
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

	t := make([]entities.Training, 0, cursor.RemainingBatchLength())
	for cursor.Next(context.Background()) {
		var training trainingData
		err = cursor.Decode(&training)
		if err != nil {
			return nil, fmt.Errorf("get user trainings: %v", err)
		}

		t = append(t, *mapTrainingToEntity(&training))
	}

	if err = cursor.Err(); err != nil {
		return nil, fmt.Errorf("get user trainings: %v", err)
	}

	return t, nil
}

func (r TrainingRepository) StartExercise(trID string, exercise *entities.TrainingExercise) (*entities.TrainingExercise, error) {
	tOID, err := primitive.ObjectIDFromHex(trID)
	if err != nil {
		return nil, errors.WithMessage(repositories.NewErrorInvalidID(trID, "training"), "start exercise")
	}

	exOID, err := primitive.ObjectIDFromHex(exercise.ExerciseID)
	if err != nil {
		return nil, repositories.NewErrorInvalidID(exercise.ExerciseID, "exercise")
	}
	newExerciseData := trainingExerciseData{
		ID:         primitive.NewObjectID(),
		ExerciseID: exOID,
		StartTime:  exercise.StartTime,
		Comment:    exercise.Comment,
		CreatedAt:  time.Now(),
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
		Sets:       mapSetsToEntities(newExerciseData.Sets),
	}
	return &newExercise, nil
}

func (r TrainingRepository) AddSet(userID, teID string, set *entities.TrainingSet) (*entities.TrainingSet, error) {
	teOID, err := primitive.ObjectIDFromHex(teID)
	if err != nil {
		return nil, errors.WithMessage(repositories.NewErrorInvalidID(teID, "training exercise"), "add set")
	}
	uOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.WithMessage(repositories.NewErrorInvalidID(userID, "user"), "add set")
	}

	newSetData := trainingSetData{
		ID:        primitive.NewObjectID(),
		Time:      set.Time,
		Reps:      set.Reps,
		CreatedAt: time.Now(),
	}

	filter := bson.M{
		"exercises._id": teOID,
		"user_id":       uOID,
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
		ID:        newSetData.ID.Hex(),
		Time:      newSetData.Time,
		Reps:      newSetData.Reps,
		CreatedAt: newSetData.CreatedAt,
	}
	return &newSet, nil
}

func (r TrainingRepository) GetTrainingExercises(id string) ([]entities.TrainingExercise, error) {
	tOID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.WithMessage(repositories.NewErrorInvalidID(id, "training"), "get training exercises")
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

	te := mapExercisesToEntities(tr[0].Exercises)
	return te, nil
}

func (r TrainingRepository) EndExercise(userID, id string, endTime time.Time) (*entities.TrainingExercise, error) {
	teOID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.WithMessage(repositories.NewErrorInvalidID(id, "training"), "end exercises")
	}
	uOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.WithMessage(repositories.NewErrorInvalidID(id, "user"), "end exercises")
	}

	filter := bson.M{"exercises._id": teOID, "user_id": uOID}
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
	te := mapExerciseToEntity(&td.Exercises[0])
	return te, nil
}
