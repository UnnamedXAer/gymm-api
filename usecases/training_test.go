package usecases_test

import (
	"testing"

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
