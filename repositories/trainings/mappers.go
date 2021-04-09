package trainings

import "github.com/unnamedxaer/gymm-api/entities"

func mapTrainingToEntity(t trainingData) entities.Training {
	return entities.Training{
		ID:        t.ID.Hex(),
		UserID:    t.UserID.Hex(),
		StartTime: t.StartTime,
		EndTime:   t.EndTime,
		Exercises: mapExercisesToEntity(t.Exercises),
	}
}

func mapExercisesToEntity(ted []trainingExerciseData) (te []entities.TrainingExercise) {

	for _, data := range ted {
		te = append(te, entities.TrainingExercise{
			ID:         data.ID.Hex(),
			ExerciseID: data.ExerciseID.Hex(),
			StartTime:  data.StartTime,
			EndTime:    data.EndTime,
			Comment:    data.Comment,
			Sets:       mapSetsToEntity(data.Sets),
		})
	}

	return te
}

func mapSetsToEntity(tsd []trainingSetData) (ts []entities.TrainingSet) {

	for _, data := range tsd {
		ts = append(ts, entities.TrainingSet{
			ID:   data.ID.Hex(),
			Time: data.Time,
			Reps: data.Reps,
		})
	}

	return ts
}
