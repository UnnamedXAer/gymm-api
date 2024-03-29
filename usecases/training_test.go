package usecases_test

import (
	"context"
	"testing"
	"time"

	"github.com/unnamedxaer/gymm-api/mocks"
	"github.com/unnamedxaer/gymm-api/usecases"
)

var (
	trainingUC    usecases.ITrainingUsecases
	trainingInput = usecases.TrainingInput{
		UserID:    mocks.ExampleTraining.UserID,
		StartTime: mocks.ExampleTraining.StartTime,
		EndTime:   mocks.ExampleTraining.EndTime,
		Comment:   mocks.ExampleTraining.Comment,
	}
)

func TestStartTraining(t *testing.T) {
	ctx := context.TODO()

	tr, err := trainingUC.StartTraining(ctx, mocks.ExampleUser.ID)
	if err != nil {
		t.Fatal(err)
	}

	if tr.StartTime.IsZero() || tr.ID == "" {
		t.Errorf("want started training, got %v", tr)
	}
}

func TestEndTraining(t *testing.T) {
	ctx := context.TODO()

	tr, err := trainingUC.EndTraining(ctx, mocks.ExampleTraining.ID)
	if err != nil {
		t.Fatal(err)
	}

	if tr.StartTime.IsZero() || tr.ID == "" {
		t.Errorf("want started training, got %v", tr)
	}
}

func TestAddTrainingSet(t *testing.T) {
	ctx := context.TODO()

	ts, err := trainingUC.AddSet(ctx, mocks.ExampleTraining.UserID, mocks.ExampleExercise.ID, &mocks.ExampleTrainingSet)
	if err != nil {
		t.Fatal(err)
	}

	if ts.ID == "" {
		t.Errorf("want set with ID, got %v", ts)
	}
}

func TestEndTrainingExercise(t *testing.T) {
	ctx := context.TODO()

	te, err := trainingUC.EndExercise(ctx, mocks.ExampleTraining.UserID, mocks.ExampleExercise.ID, time.Now())
	if err != nil {
		t.Fatal(err)
	}

	if te.EndTime.IsZero() {
		t.Errorf("want ended exercise, got %v", te)
	}
}

func TestAddTrainingExercise(t *testing.T) {
	ctx := context.TODO()

	tr, err := trainingUC.EndTraining(ctx, mocks.ExampleTraining.ID)
	if err != nil {
		t.Fatal(err)
	}

	if tr.StartTime.IsZero() || tr.ID == "" {
		t.Errorf("want started exercise, got %v", tr)
	}
}

func TestGetTrainingExercises(t *testing.T) {
	ctx := context.TODO()

	te, err := trainingUC.GetTrainingExercises(ctx, mocks.ExampleTraining.ID)
	if err != nil {
		t.Fatal(err)
	}

	if len(te) != len(mocks.ExampleTraining.Exercises) {
		t.Errorf("want %d exercises, got %d", len(te), len(mocks.ExampleTraining.Exercises))
	}
}

func TestGetUserTrainings(t *testing.T) {
	ctx := context.TODO()

	tr, err := trainingUC.GetUserTrainings(ctx, mocks.ExampleTraining.UserID, true)
	if err != nil {
		t.Fatal(err)
	}

	cnt := len(tr)
	if cnt == 0 {
		t.Errorf("want not empty slice of trainings, got %v", tr)
		return
	}

	for _, v := range tr {
		if !v.EndTime.IsZero() {
			t.Errorf("want only started trainings, got %v", tr)
			return
		}
	}

	tr, err = trainingUC.GetUserTrainings(ctx, mocks.ExampleTraining.UserID, false)
	if err != nil {
		t.Fatal(err)
	}

	if len(tr) < cnt {
		t.Errorf("want at least %d trainings, got %d", cnt, len(tr))
	}
}
