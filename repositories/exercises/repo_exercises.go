package exercises

import (
	"context"
	"fmt"
	"time"

	"github.com/unnamedxaer/gymm-api/entities"
	"github.com/unnamedxaer/gymm-api/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ExerciseData struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name,omitempty"`
	Description string             `bson:"description,omitempty"`
	SetUnit     entities.SetUnit   `bson:"set_unit,omitempty"`
	CreatedAt   time.Time          `bson:"created_at,omitempty"`
	CreatedBy   string             `bson:"created_by,omitempty"`
}

func (repo ExerciseRepository) GetExerciseByID(id string) (*entities.Exercise, error) {
	exOID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("get exercise by id: %v", err)
	}

	filter := bson.M{"_id": exOID}

	result := repo.col.FindOne(context.Background(), filter)
	if err = result.Err(); err != nil {
		return nil, fmt.Errorf("get exercise by id: %v", err)
	}

	var data ExerciseData
	err = result.Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("get exercise by id: %v", err)
	}
	ex := mapExerciseToEntity(&data)
	return ex, nil
}

func (repo ExerciseRepository) CreateExercise(name, description string, setUnit entities.SetUnit, createdBy string) (*entities.Exercise, error) {

	data := ExerciseData{
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		CreatedBy:   createdBy,
		SetUnit:     setUnit,
	}

	result, err := repo.col.InsertOne(context.Background(), &data)
	if err != nil {
		return nil, fmt.Errorf("create exercise: %v", err)
	}

	var ok bool
	data.ID, ok = result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("create exercise: %v", fmt.Errorf("type assertion: InsertedID is not of type primitive.ObjectID"))
	}

	ex := mapExerciseToEntity(&data)

	return ex, nil
}

func (repo ExerciseRepository) UpdateExercise(ex *entities.Exercise) (*entities.Exercise, error) {
	exOID, err := primitive.ObjectIDFromHex(ex.ID)
	if err != nil {
		return nil, repositories.NewErrorInvalidID(ex.ID)
	}

	filter := bson.M{
		"_id": exOID,
	}

	update := bson.D{}
	if ex.Name != "" {
		update = append(update, primitive.E{"name", ex.Name})
	}
	if ex.Description != "" {
		update = append(update, primitive.E{"description", ex.Description})
	}
	if ex.SetUnit != 0 {
		update = append(update, primitive.E{"set_unit", ex.SetUnit})
	}

	update = bson.D{{"$set", update}}

	result := repo.col.FindOneAndUpdate(context.Background(), filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After))
	if err = result.Err(); err != nil {
		return nil, fmt.Errorf("update exercise: %v", err)
	}
	var data ExerciseData
	err = result.Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("update exercise: %v", err)
	}

	updatedEx := mapExerciseToEntity(&data)

	return updatedEx, nil
}
