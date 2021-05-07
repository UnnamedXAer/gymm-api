package usecases

import (
	"context"
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
	// GetTrainingByID returns training for given id
	GetTrainingByID(ctx context.Context, id string) (*entities.Training, error)
	// StartTraining starts new training by inserting new record in training storage with start time.
	StartTraining(ctx context.Context, userID string, startTime time.Time) (*entities.Training, error)
	// EndTraining marks given training as completed by setting training end time.
	EndTraining(ctx context.Context, trainingID string, endTime time.Time) (*entities.Training, error)
	GetUserTrainings(ctx context.Context, userID string, started bool) (t []entities.Training, err error)
	StartExercise(ctx context.Context, trID string, exercise *entities.TrainingExercise) (*entities.TrainingExercise, error)
	AddSet(ctx context.Context, userID, teID string, set *entities.TrainingSet) (*entities.TrainingSet, error)
	GetTrainingExercises(ctx context.Context, id string) ([]entities.TrainingExercise, error)
	EndExercise(ctx context.Context, userID, id string, endTime time.Time) (*entities.TrainingExercise, error)
}

type TrainingUsecases struct {
	repo TrainingRepo
}

type ITrainingUsecases interface {
	GetTrainingByID(ctx context.Context, id string) (*entities.Training, error)
	StartTraining(ctx context.Context, userID string) (*entities.Training, error)
	EndTraining(ctx context.Context, id string) (*entities.Training, error)
	GetUserTrainings(ctx context.Context, userID string, started bool) (t []entities.Training, err error)
	StartExercise(ctx context.Context, trID string, exercise *entities.TrainingExercise) (*entities.TrainingExercise, error)
	AddSet(ctx context.Context, userID, teID string, set *entities.TrainingSet) (*entities.TrainingSet, error)
	GetTrainingExercises(ctx context.Context, id string) ([]entities.TrainingExercise, error)
	EndExercise(ctx context.Context, userID, id string, endTime time.Time) (*entities.TrainingExercise, error)
}

// GetTrainingByID returns training for given id
func (tu *TrainingUsecases) GetTrainingByID(ctx context.Context,
	id string) (*entities.Training, error) {
	return tu.repo.GetTrainingByID(ctx, id)
}

// StartTraining creates a new training.
func (tu *TrainingUsecases) StartTraining(ctx context.Context,
	userID string) (*entities.Training, error) {
	return tu.repo.StartTraining(ctx, userID, time.Now())
}

// EndTraining stops current training.
func (tu *TrainingUsecases) EndTraining(ctx context.Context,
	id string) (*entities.Training, error) {
	return tu.repo.EndTraining(ctx, id, time.Now())
}

func (tu *TrainingUsecases) GetUserTrainings(ctx context.Context,
	userID string, started bool) (t []entities.Training, err error) {
	return tu.repo.GetUserTrainings(ctx, userID, started)
}

func (tu *TrainingUsecases) StartExercise(ctx context.Context,
	trID string, exercise *entities.TrainingExercise) (*entities.TrainingExercise, error) {
	return tu.repo.StartExercise(ctx, trID, exercise)
}

func (tu *TrainingUsecases) AddSet(ctx context.Context,
	userID, teID string, set *entities.TrainingSet) (*entities.TrainingSet, error) {
	return tu.repo.AddSet(ctx, userID, teID, set)
}

func (tu *TrainingUsecases) GetTrainingExercises(ctx context.Context,
	id string) ([]entities.TrainingExercise, error) {
	return tu.repo.GetTrainingExercises(ctx, id)
}

func (tu *TrainingUsecases) EndExercise(ctx context.Context,
	userID, id string, endTime time.Time) (*entities.TrainingExercise, error) {
	return tu.repo.EndExercise(ctx, userID, id, endTime)
}

func NewTrainingUseCases(repo TrainingRepo) ITrainingUsecases {
	return &TrainingUsecases{
		repo: repo,
	}
}
