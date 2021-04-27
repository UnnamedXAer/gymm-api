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

func mapExercisesToEntity(exd []ExerciseData) (exercises []entities.Exercise) {

	for _, data := range exd {
		exercises = append(exercises, mapExerciseToEntity(&data))
	}

	return exercises
}
