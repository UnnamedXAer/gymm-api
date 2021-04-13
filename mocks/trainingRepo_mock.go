package mocks

import (
	"time"

	"github.com/unnamedxaer/gymm-api/entities"
)

var (
	ExampleTraining = entities.Training{
		ID:        "607443ceb40d9ea8602803e7",
		UserID:    UserID,
		StartTime: time.Now().Add(2 * time.Hour),
		Exercises: []entities.TrainingExercise{
			{
				ID:         "607400d5bf81935a539bd698",
				ExerciseID: ExampleExercise.ID,
				StartTime:  time.Now().Add(2 * time.Hour),
				EndTime:    time.Now(),
			},
		},
		EndTime: time.Now(),
		Comment: "too long, too heavy",
	}
)

type MockTrainingRepo struct {
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
