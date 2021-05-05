package mocks

import (
	"strings"
	"time"

	"github.com/unnamedxaer/gymm-api/entities"
	"github.com/unnamedxaer/gymm-api/repositories"
	"github.com/unnamedxaer/gymm-api/usecases"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ExampleExercise = entities.Exercise{
	ID:          "6072d3206144644984a54fa0",
	Name:        "Deadlift",
	Description: "The deadlift is an exercise in which a loaded bar is lifted off the ground to the level of the hips.",
	SetUnit:     entities.Weight,
	CreatedAt:   time.Now().UTC(),
	CreatedBy:   UserID,
}

func InsertMockExercise(er usecases.ExerciseRepo) (*entities.Exercise, error) {

	return er.CreateExercise(
		ExampleExercise.Name,
		ExampleExercise.Description,
		ExampleExercise.SetUnit,
		ExampleExercise.CreatedBy,
	)
}

type MockExerciseRepo struct{}

func (er *MockExerciseRepo) CreateExercise(name, description string, setUnit entities.SetUnit, createdBy string) (*entities.Exercise, error) {

	return &entities.Exercise{
		ID:          ExampleExercise.ID,
		Name:        name,
		Description: description,
		SetUnit:     setUnit,
		CreatedAt:   ExampleExercise.CreatedAt,
		CreatedBy:   createdBy,
	}, nil
}

func (er *MockExerciseRepo) GetExerciseByID(id string) (*entities.Exercise, error) {
	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, repositories.NewErrorInvalidID(id, "exercise")
	}

	if ExampleExercise.ID == id {
		out := ExampleExercise
		return &out, nil
	}

	return nil, nil //repositories.NewErrorNotFoundRecord()
}

func (er *MockExerciseRepo) GetExercisesByName(name string) ([]entities.Exercise, error) {

	if strings.Contains(strings.ToLower(ExampleExercise.Name), strings.ToLower(name)) {
		out := []entities.Exercise{ExampleExercise}
		return out, nil
	}

	return nil, nil
}

func (er *MockExerciseRepo) UpdateExercise(ex *entities.Exercise) (*entities.Exercise, error) {
	_, err := primitive.ObjectIDFromHex(ex.ID)
	if err != nil {
		return nil, repositories.NewErrorInvalidID(ex.ID, "exercise")
	}

	out := ExampleExercise
	out.ID = ex.ID
	if ex.Description != "" {
		out.Description = ex.Description
	}
	if ex.Name != "" {
		out.Name = ex.Name
	}
	if ex.SetUnit != 0 {
		out.SetUnit = ex.SetUnit
	}

	return &out, nil
}
