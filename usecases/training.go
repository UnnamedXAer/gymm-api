package usecases

import "github.com/unnamedxaer/gymm-api/entities"

// TrainingRepo represents trainings repository
type TrainingRepo interface {
	// StartTraining starts new training by inserting new records in training storage with start time
	StartTraining() (entities.Training, error)
	// EndTraining marks given training as completed.
	EndTraining() (entities.Training, error)
}

type TrainingUsecases struct {
	repo TrainingRepo
}

type ITrainingUsecases interface {
	StartTraining()
	EndTraining()
}

// StartTraining creates a new training. If there is not stopped training for current user the new training will not be created.
func StartTraining() {
	panic("not implemented yet")
}

// EndTraining stops current training.
func EndTraining() {
	panic("not implemented yet")
}
