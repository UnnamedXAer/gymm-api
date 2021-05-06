package exercises

import "github.com/unnamedxaer/gymm-api/entities"

func mapExerciseToEntity(data *ExerciseData) entities.Exercise {
	return entities.Exercise{
		ID:          data.ID.Hex(),
		Name:        data.Name,
		Description: data.Description,
		SetUnit:     data.SetUnit,
		CreatedAt:   data.CreatedAt.UTC(),
		CreatedBy:   data.CreatedBy,
	}
}

func mapExercisesToEntities(exd []ExerciseData) []entities.Exercise {

	exercises := make([]entities.Exercise, len(exd))

	for i := 0; i < len(exd); i++ {
		exercises[i] = mapExerciseToEntity(&exd[i])
	}

	return exercises
}
