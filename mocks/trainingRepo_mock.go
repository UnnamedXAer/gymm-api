package mocks

import (
	"strings"
	"time"

	"github.com/unnamedxaer/gymm-api/entities"
	"github.com/unnamedxaer/gymm-api/repositories"
)

var (
	ExampleTraining = entities.Training{
		ID:        "607443ceb40d9ea8602803e7",
		UserID:    UserID,
		StartTime: time.Now().Add(-2 * time.Hour),
		Exercises: []entities.TrainingExercise{
			{
				ID:         "607400d5bf81935a539bd698",
				ExerciseID: ExampleExercise.ID,
				StartTime:  time.Now().Add(-115 * time.Minute),
				EndTime:    time.Now(),
				Sets: []entities.TrainingSet{
					{
						ID:   "60740f289ee8e963adb5412a",
						Time: time.Now().Add(-110 * time.Minute),
						Reps: 12,
					},
					{
						ID:   "60740f289ee8e963adb5412d",
						Time: time.Now().Add(-107 * time.Minute),
						Reps: 10,
					},
					{
						ID:   "60740f289ee8e963adb5412c",
						Time: time.Now().Add(-103 * time.Minute),
						Reps: 10,
					},
				},
				Comment: "too short breaks",
			},
			{
				ID:         "60740f289ee8e963adb5412b",
				ExerciseID: ExampleExercise.ID,
				StartTime:  time.Now().Add(-115 * time.Minute),
			},
		},
		EndTime: time.Now(),
		Comment: "too long, too heavy",
	}

	ExampleTrainingExercise = ExampleTraining.Exercises[0]
	ExampleTrainingSet      = ExampleTrainingExercise.Sets[0]
)

type MockTrainingRepo struct {
}

func (tr *MockTrainingRepo) GetTrainingByID(id string) (*entities.Training, error) {

	if strings.Contains(id, "notfound") {
		return nil, nil
	}

	if strings.Contains(id, "INVALIDID") {
		return nil, repositories.NewErrorInvalidID(id, "training")
	}

	out := ExampleTraining
	out.ID = id
	return &out, nil
}

func (tr *MockTrainingRepo) StartTraining(userID string, startTime time.Time) (*entities.Training, error) {
	return &entities.Training{
		ID:        userID,
		UserID:    ExampleTraining.UserID,
		StartTime: ExampleTraining.StartTime,
	}, nil
}

func (tr *MockTrainingRepo) EndTraining(id string, endTime time.Time) (*entities.Training, error) {
	out := ExampleTraining
	out.ID = id
	return &out, nil
}

func (tr *MockTrainingRepo) GetUserTrainings(userID string, started bool) ([]entities.Training, error) {
	out := []entities.Training{ExampleTraining}

	out[0].UserID = userID
	if !started {
		out[0].EndTime = time.Time{}
	}
	return out, nil
}

func (tr *MockTrainingRepo) StartExercise(trID string, exercise *entities.TrainingExercise) (*entities.TrainingExercise, error) {
	out := ExampleTrainingExercise
	out.EndTime = time.Time{}
	return &out, nil
}

func (tr *MockTrainingRepo) AddSet(userID, teID string, set *entities.TrainingSet) (*entities.TrainingSet, error) {
	out := *set
	out.ID = ExampleTrainingSet.ID
	out.Time = ExampleTrainingSet.Time
	return &out, nil
}

func (tr *MockTrainingRepo) GetTrainingExercises(id string) ([]entities.TrainingExercise, error) {
	out := []entities.TrainingExercise{}
	for i, ex := range ExampleTraining.Exercises {
		out = append(out, ex)
		out[i].ID = id
	}
	return out, nil
}

func (tr *MockTrainingRepo) EndExercise(userID, id string, endTime time.Time) (*entities.TrainingExercise, error) {
	out := ExampleTrainingExercise
	out.ID = id
	out.EndTime = endTime
	return &out, nil
}
