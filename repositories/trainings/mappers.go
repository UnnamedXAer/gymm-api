package trainings

import "github.com/unnamedxaer/gymm-api/entities"

func mapTrainingToEntity(td trainingData) entities.Training {
	return entities.Training{
		ID:        td.ID.Hex(),
		UserID:    td.UserID.Hex(),
		StartTime: td.StartTime,
		EndTime:   td.EndTime,
		Exercises: mapExercisesToEntity(td.Exercises),
	}
}

func mapExerciseToEntity(ted trainingExerciseData) (te entities.TrainingExercise) {
	return entities.TrainingExercise{
		ID:         ted.ID.Hex(),
		ExerciseID: ted.ExerciseID.Hex(),
		StartTime:  ted.StartTime,
		EndTime:    ted.EndTime,
		Comment:    ted.Comment,
		Sets:       mapSetsToEntity(ted.Sets),
	}
}

func mapExercisesToEntity(ted []trainingExerciseData) (te []entities.TrainingExercise) {

	for _, data := range ted {
		te = append(te, mapExerciseToEntity(data))
	}

	return te
}

func mapSetToEntity(tsd trainingSetData) (ts entities.TrainingSet) {
	return entities.TrainingSet{
		ID:   tsd.ID.Hex(),
		Time: tsd.Time,
		Reps: tsd.Reps,
	}
}

func mapSetsToEntity(tsd []trainingSetData) (ts []entities.TrainingSet) {

	for _, data := range tsd {
		ts = append(ts, mapSetToEntity(data))
	}

	return ts
}
