package usecases

import (
	"time"

	"github.com/unnamedxaer/gymm-api/entities"
)

type ExerciseInput struct {
	Name        string           `json:"name" validate:"required,min=2,max=50"`
	Description string           `json:"description" validate:"required,min=20,max=500"`
	SetUnit     entities.SetUnit `json:"setUnit" validate:"required"`
	CreatedAt   time.Time        `json:"createdAt" validate:""`
	CreatedBy   string           `json:"createdBy" validate:"required"`
}

type ExerciseRepo interface {
	CreateExercise(name, description string, setUnit entities.SetUnit, createdBy string) (*entities.Exercise, error)
	GetExerciseByID(id string) (*entities.Exercise, error)
	UpdateExercise(ex *entities.Exercise) (*entities.Exercise, error)
}

type ExerciseUseCases struct {
	repo ExerciseRepo
}

type IExerciseUseCases interface {
	CreateExercise(name, description string, setUnit entities.SetUnit, loggedUserID string) (*entities.Exercise, error)
	GetExerciseByID(id string) (*entities.Exercise, error)
	UpdateExercise(ex *entities.Exercise) (*entities.Exercise, error)
}

func (eu *ExerciseUseCases) CreateExercise(name, description string, setUnit entities.SetUnit, loggedUserID string) (*entities.Exercise, error) {
	return eu.repo.CreateExercise(name, description, setUnit, loggedUserID)
}

func (eu *ExerciseUseCases) GetExerciseByID(id string) (*entities.Exercise, error) {
	return eu.repo.GetExerciseByID(id)
}

func (eu *ExerciseUseCases) UpdateExercise(ex *entities.Exercise) (*entities.Exercise, error) {
	return eu.repo.UpdateExercise(ex)
}

func NewExerciseUseCases(exRepo ExerciseRepo) IExerciseUseCases {
	return &ExerciseUseCases{
		repo: exRepo,
	}
}
