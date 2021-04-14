package mocks

import (
	"time"

	"github.com/unnamedxaer/gymm-api/entities"
	"github.com/unnamedxaer/gymm-api/repositories"
	"github.com/unnamedxaer/gymm-api/usecases"
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
	if ExampleExercise.ID == id {
		out := ExampleExercise
		return &out, nil
	}

	return nil, repositories.NewErrorNotFoundRecord()
}

func (er *MockExerciseRepo) UpdateExercise(ex *entities.Exercise) (*entities.Exercise, error) {
	if ex.ID == "" {
		return nil, repositories.NewErrorInvalidID(ex.ID)
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
