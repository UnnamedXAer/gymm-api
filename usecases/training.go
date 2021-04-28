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
	// GetTrainingByID returns training for given id
	GetTrainingByID(id string) (*entities.Training, error)
	// StartTraining starts new training by inserting new record in training storage with start time.
	StartTraining(userID string, startTime time.Time) (*entities.Training, error)
	// EndTraining marks given training as completed by setting training end time.
	EndTraining(trainingID string, endTime time.Time) (*entities.Training, error)
	GetUserTrainings(userID string, started bool) (t []entities.Training, err error)
	AddExercise(trID string, exercise *entities.TrainingExercise) (*entities.TrainingExercise, error)
	AddSet(teID string, set *entities.TrainingSet) (*entities.TrainingSet, error)
	GetTrainingExercises(id string) ([]entities.TrainingExercise, error)
	EndExercise(id string, endTime time.Time) (*entities.TrainingExercise, error)
}

type TrainingUsecases struct {
	repo TrainingRepo
}

type ITrainingUsecases interface {
	GetTrainingByID(id string) (*entities.Training, error)
	StartTraining(userID string) (*entities.Training, error)
	EndTraining(id string) (*entities.Training, error)
	GetUserTrainings(userID string, started bool) (t []entities.Training, err error)
	AddExercise(trID string, exercise *entities.TrainingExercise) (*entities.TrainingExercise, error)
	AddSet(teID string, set *entities.TrainingSet) (*entities.TrainingSet, error)
	GetTrainingExercises(id string) ([]entities.TrainingExercise, error)
	EndExercise(id string, endTime time.Time) (*entities.TrainingExercise, error)
}

// GetTrainingByID returns training for given id
func (tu *TrainingUsecases) GetTrainingByID(id string) (*entities.Training, error) {
	return tu.repo.GetTrainingByID(id)
}

// StartTraining creates a new training.
func (tu *TrainingUsecases) StartTraining(userID string) (*entities.Training, error) {
	return tu.repo.StartTraining(userID, time.Now())
}

// EndTraining stops current training.
func (tu *TrainingUsecases) EndTraining(id string) (*entities.Training, error) {
	return tu.repo.EndTraining(id, time.Now())
}

func (tu *TrainingUsecases) GetUserTrainings(userID string, started bool) (t []entities.Training, err error) {
	return tu.repo.GetUserTrainings(userID, started)
}

func (tu *TrainingUsecases) AddExercise(trID string, exercise *entities.TrainingExercise) (*entities.TrainingExercise, error) {
	return tu.repo.AddExercise(trID, exercise)
}

func (tu *TrainingUsecases) AddSet(teID string, set *entities.TrainingSet) (*entities.TrainingSet, error) {
	return tu.repo.AddSet(teID, set)
}

func (tu *TrainingUsecases) GetTrainingExercises(id string) ([]entities.TrainingExercise, error) {
	return tu.repo.GetTrainingExercises(id)
}

func (tu *TrainingUsecases) EndExercise(id string, endTime time.Time) (*entities.TrainingExercise, error) {
	return tu.repo.EndExercise(id, endTime)
}

func NewTrainingUseCases(repo TrainingRepo) ITrainingUsecases {
	return &TrainingUsecases{
		repo: repo,
	}
}
