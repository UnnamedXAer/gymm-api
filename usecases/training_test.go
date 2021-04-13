package usecases_test

import (
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
	tr, err := trainingUC.StartTraining(mocks.ExampleUser.ID)
	if err != nil {
		t.Fatal(err)
	}

	if tr.StartTime.IsZero() || tr.ID == "" {
		t.Errorf("want started training, got %v", tr)
	}
}

func TestEndTraining(t *testing.T) {
	tr, err := trainingUC.EndTraining(mocks.ExampleTraining.ID)
	if err != nil {
		t.Fatal(err)
	}

	if tr.StartTime.IsZero() || tr.ID == "" {
		t.Errorf("want started training, got %v", tr)
	}
}

func TestAddTrainingSet(t *testing.T) {
	ts, err := trainingUC.AddSet(mocks.ExampleExercise.ID, &mocks.ExampleTrainingSet)
	if err != nil {
		t.Fatal(err)
	}

	if ts.ID == "" {
		t.Errorf("want set with ID, got %v", ts)
	}
}

func TestEndTrainingExercise(t *testing.T) {
	te, err := trainingUC.EndExercise(mocks.ExampleExercise.ID, time.Now())
	if err != nil {
		t.Fatal(err)
	}

	if te.EndTime.IsZero() {
		t.Errorf("want ended exercise, got %v", te)
	}
}

func TestAddTreiningExercise(t *testing.T) {
	tr, err := trainingUC.EndTraining(mocks.ExampleTraining.ID)
	if err != nil {
		t.Fatal(err)
	}

	if tr.StartTime.IsZero() || tr.ID == "" {
		t.Errorf("want started exercise, got %v", tr)
	}
}

func TestGetTrainingExercises(t *testing.T) {
	te, err := trainingUC.GetTrainingExercises(mocks.ExampleTraining.ID)
	if err != nil {
		t.Fatal(err)
	}

	if len(te) != len(mocks.ExampleTraining.Exercises) {
		t.Errorf("want %d exercises, got %d", len(te), len(mocks.ExampleTraining.Exercises))
	}
}

func TestGetUserTrainings(t *testing.T) {
	tr, err := trainingUC.GetUserTrainings(mocks.ExampleTraining.UserID, true)
	if err != nil {
		t.Fatal(err)
	}

	cnt := len(tr)
	if cnt == 0 {
		t.Errorf("want not empty slice of trainings, got %v", tr)
		return
	}

	for _, v := range tr {
		if v.EndTime.IsZero() {
			t.Errorf("want only started trainings, got %v", tr)
			return
		}
	}

	tr, err = trainingUC.GetUserTrainings(mocks.ExampleTraining.UserID, false)
	if err != nil {
		t.Fatal(err)
	}

	if len(tr) < cnt {
		t.Errorf("want at least %d trainings, got %d", cnt, len(tr))
	}
}
