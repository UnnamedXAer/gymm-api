package exercises

import (
	"context"
	"fmt"
	"time"

	"github.com/unnamedxaer/gymm-api/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	var ex entities.Exercise
	err = result.Decode(&ex)
	if err != nil {
		return nil, fmt.Errorf("get exercise by id: %v", err)
	}
	return &ex, nil
}

func (repo ExerciseRepository) CreateExercise(name, description string, setUnit entities.SetUnit) (*entities.Exercise, error) {
	panic("not implemented yet")
}

func (repo ExerciseRepository) UpdateExercise(id string, name, description string, setUnit uint8) (*entities.Exercise, error) {
	panic("not implemented yet")
}
