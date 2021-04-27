package usecases

import (
	"time"

	"github.com/unnamedxaer/gymm-api/entities"
)

type ExerciseInput struct {
	Name        string           `json:"name" validate:"required,min=2,max=50,ex_name_chars,printascii"`
	Description string           `json:"description" validate:"required,min=10,max=500,printascii"`
	SetUnit     entities.SetUnit `json:"setUnit" validate:"set_unit,required,oneof=1 2"`
	CreatedAt   time.Time        `json:"createdAt" validate:"-"`
	CreatedBy   string           `json:"createdBy" validate:"-"`
}

type ExerciseRepo interface {
	CreateExercise(name, description string, setUnit entities.SetUnit, createdBy string) (*entities.Exercise, error)
	GetExerciseByID(id string) (*entities.Exercise, error)
	GetExercisesByName(name string) ([]entities.Exercise, error)
	UpdateExercise(ex *entities.Exercise) (*entities.Exercise, error)
}

type ExerciseUseCases struct {
	repo ExerciseRepo
}

type IExerciseUseCases interface {
	CreateExercise(name, description string, setUnit entities.SetUnit, loggedUserID string) (*entities.Exercise, error)
	GetExerciseByID(id string) (*entities.Exercise, error)
	GetExercisesByName(name string) ([]entities.Exercise, error)
	UpdateExercise(ex *entities.Exercise) (*entities.Exercise, error)
}

func (eu *ExerciseUseCases) CreateExercise(name, description string, setUnit entities.SetUnit, loggedUserID string) (*entities.Exercise, error) {
	return eu.repo.CreateExercise(name, description, setUnit, loggedUserID)
}

func (eu *ExerciseUseCases) GetExerciseByID(id string) (*entities.Exercise, error) {
	return eu.repo.GetExerciseByID(id)
}
func (eu *ExerciseUseCases) GetExercisesByName(name string) ([]entities.Exercise, error) {
	return eu.repo.GetExercisesByName(name)
}

func (eu *ExerciseUseCases) UpdateExercise(ex *entities.Exercise) (*entities.Exercise, error) {
	return eu.repo.UpdateExercise(ex)
}

func NewExerciseUseCases(exRepo ExerciseRepo) IExerciseUseCases {
	return &ExerciseUseCases{
		repo: exRepo,
	}
}
