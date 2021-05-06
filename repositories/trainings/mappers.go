package trainings

import "github.com/unnamedxaer/gymm-api/entities"

func mapTrainingToEntity(td *trainingData) *entities.Training {
	return &entities.Training{
		ID:        td.ID.Hex(),
		UserID:    td.UserID.Hex(),
		StartTime: td.StartTime,
		EndTime:   td.EndTime,
		Exercises: mapExercisesToEntities(td.Exercises),
		Comment:   td.Comment,
		CreatedAt: td.CreatedAt,
	}
}

func mapExerciseToEntity(ted *trainingExerciseData) *entities.TrainingExercise {
	return &entities.TrainingExercise{
		ID:         ted.ID.Hex(),
		ExerciseID: ted.ExerciseID.Hex(),
		StartTime:  ted.StartTime,
		EndTime:    ted.EndTime,
		Comment:    ted.Comment,
		Sets:       mapSetsToEntities(ted.Sets),
		CreatedAt:  ted.CreatedAt,
	}
}

func mapExercisesToEntities(ted []trainingExerciseData) []entities.TrainingExercise {

	te := make([]entities.TrainingExercise, len(ted))

	for i := 0; i < len(ted); i++ {
		te[i] = *mapExerciseToEntity(&ted[i])
	}

	return te
}

func mapSetToEntity(tsd trainingSetData) *entities.TrainingSet {
	return &entities.TrainingSet{
		ID:        tsd.ID.Hex(),
		Time:      tsd.Time,
		Reps:      tsd.Reps,
		CreatedAt: tsd.CreatedAt,
	}
}

func mapSetsToEntities(tsd []trainingSetData) []entities.TrainingSet {

	ts := make([]entities.TrainingSet, len(tsd))

	for i := 0; i < len(tsd); i++ {
		ts[i] = *mapSetToEntity(tsd[i])
	}

	return ts
}
