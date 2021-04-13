package usecases

import (
	"time"

	"github.com/unnamedxaer/gymm-api/entities"
)

type TrainingInput struct {
	UserID    string    `json:"userId"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	// Exercises []ExerciseInput `json:"exercises"`
	Comment string `json:"comment"`
}

// TrainingRepo represents trainings repository
type TrainingRepo interface {
	// StartTraining starts new training by inserting new record in training storage with start time.
	StartTraining(userID string, startTime time.Time) (*entities.Training, error)
	// EndTraining marks given training as completed by setting training end time.
	EndTraining(trainingID string, endTime time.Time) (*entities.Training, error)
}

type TrainingUsecases struct {
	repo TrainingRepo
}

type ITrainingUsecases interface {
	StartTraining(userID string) (*entities.Training, error)
	EndTraining(id string) (*entities.Training, error)
}

// StartTraining creates a new training.
func (tu *TrainingUsecases) StartTraining(userID string) (*entities.Training, error) {
	return tu.repo.StartTraining(userID, time.Now())
}

// EndTraining stops current training.
func (tu *TrainingUsecases) EndTraining(ID string) (*entities.Training, error) {
	return tu.repo.EndTraining(ID, time.Now())
}

func NewTrainingUseCases(repo TrainingRepo) ITrainingUsecases {
	return &TrainingUsecases{
		repo: repo,
	}
}
