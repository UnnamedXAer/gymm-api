package usecases

import "github.com/unnamedxaer/gymm-api/entities"

type ExerciseRepo interface {
	CreateExercise(name, description string, setUnit entities.SetUnit) (*entities.Exercise, error)
	GetExerciseByID(id string) (*entities.Exercise, error)
	UpdateExercise(id string, name, description string, setUnit uint8) (*entities.Exercise, error)
}

type ExerciseUseCases struct {
	repo ExerciseRepo
}

type IExerciseUseCases interface {
	CreateExercise(name, description string, setUnit uint8) (*entities.Exercise, error)
	GetExerciseByID(id string) (*entities.Exercise, error)
	UpdateExercise(id string, name, description string, setUnit uint8) (*entities.Exercise, error)
}

func (eu *ExerciseUseCases) CreateExercise(name, description string, setUnit uint8) (*entities.Exercise, error) {
	panic("not implemented yet")
}

func (eu *ExerciseUseCases) GetExerciseByID(id string) (*entities.Exercise, error) {
	panic("not implemented yet")
}

func (eu *ExerciseUseCases) UpdateExercise(id string, name, description string, setUnit uint8) (*entities.Exercise, error) {
	panic("not implemented yet")
}

func NewExerciseUseCases(exRepo ExerciseRepo) IExerciseUseCases {
	return &ExerciseUseCases{
		repo: exRepo,
	}
}
