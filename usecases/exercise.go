package usecases

import (
	"context"
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
	CreateExercise(ctx context.Context, name, description string, setUnit entities.SetUnit, createdBy string) (*entities.Exercise, error)
	GetExerciseByID(ctx context.Context, id string) (*entities.Exercise, error)
	GetExercisesByName(ctx context.Context, name string) ([]entities.Exercise, error)
	UpdateExercise(ctx context.Context, ex *entities.Exercise) (*entities.Exercise, error)
}

type ExerciseUseCases struct {
	repo ExerciseRepo
}

type IExerciseUseCases interface {
	CreateExercise(ctx context.Context, name, description string, setUnit entities.SetUnit, loggedUserID string) (*entities.Exercise, error)
	GetExerciseByID(ctx context.Context, id string) (*entities.Exercise, error)
	GetExercisesByName(ctx context.Context, name string) ([]entities.Exercise, error)
	UpdateExercise(ctx context.Context, ex *entities.Exercise) (*entities.Exercise, error)
}

func (eu *ExerciseUseCases) CreateExercise(
	ctx context.Context,
	name string,
	description string,
	setUnit entities.SetUnit,
	loggedUserID string) (*entities.Exercise, error) {
	return eu.repo.CreateExercise(ctx, name, description, setUnit, loggedUserID)
}

func (eu *ExerciseUseCases) GetExerciseByID(
	ctx context.Context,
	id string) (*entities.Exercise, error) {
	return eu.repo.GetExerciseByID(ctx, id)
}
func (eu *ExerciseUseCases) GetExercisesByName(
	ctx context.Context,
	name string) ([]entities.Exercise, error) {
	return eu.repo.GetExercisesByName(ctx, name)
}

func (eu *ExerciseUseCases) UpdateExercise(
	ctx context.Context,
	ex *entities.Exercise) (*entities.Exercise, error) {
	return eu.repo.UpdateExercise(ctx, ex)
}

func NewExerciseUseCases(exRepo ExerciseRepo) IExerciseUseCases {
	return &ExerciseUseCases{
		repo: exRepo,
	}
}
