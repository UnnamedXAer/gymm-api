package usecases_test

import (
	"context"
	"testing"
	"time"

	"github.com/unnamedxaer/gymm-api/entities"
	"github.com/unnamedxaer/gymm-api/mocks"
	"github.com/unnamedxaer/gymm-api/testhelpers"
	"github.com/unnamedxaer/gymm-api/usecases"
)

var (
	exerciseUC    usecases.IExerciseUseCases
	exerciseInput = usecases.ExerciseInput{
		Name:        mocks.ExampleExercise.Name,
		Description: mocks.ExampleExercise.Description,
		SetUnit:     mocks.ExampleExercise.SetUnit,
		CreatedBy:   mocks.ExampleExercise.CreatedBy,
	}
)

func TestCreateExercise(t *testing.T) {
	ctx := context.TODO()

	got, _ := exerciseUC.CreateExercise(ctx, exerciseInput.Name, exerciseInput.Description, exerciseInput.SetUnit, exerciseInput.CreatedBy)
	if got.ID == "" ||
		got.Name != exerciseInput.Name ||
		got.Description != exerciseInput.Description ||
		got.CreatedAt.IsZero() ||
		got.SetUnit != exerciseInput.SetUnit ||
		got.CreatedBy != exerciseInput.CreatedBy {
		t.Fatalf("want %v got %v", exerciseInput, got)
	}
}

func TestGetExerciseByID(t *testing.T) {
	ctx := context.TODO()

	got, _ := exerciseUC.GetExerciseByID(ctx, mocks.ExampleExercise.ID)
	if got.ID != mocks.ExampleExercise.ID {
		t.Fatalf("want\n%v got\n%v", mocks.ExampleExercise, got)
	}
}

func TestUpdateExercise(t *testing.T) {
	ctx := context.TODO()

	var input entities.Exercise
	input.ID = mocks.ExampleExercise.ID
	input.Description = mocks.ExampleExercise.Description + "\n->" + time.Now().String()
	got, _ := exerciseUC.UpdateExercise(ctx, &input)
	if got.ID != mocks.ExampleExercise.ID {
		t.Fatalf("want\n%v got\n%v", mocks.ExampleExercise, got)
	}

	if got.Name != exerciseInput.Name ||
		got.Description != input.Description ||
		!testhelpers.TimesEqual(
			got.CreatedAt, mocks.ExampleExercise.CreatedAt) ||
		got.SetUnit != exerciseInput.SetUnit ||
		got.CreatedBy != exerciseInput.CreatedBy {
		t.Fatalf("want\n%v got\n%v", exerciseInput, got)
	}
}
