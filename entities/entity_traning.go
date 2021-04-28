package entities

import "time"

// Training keeps an informations about set of executed exercises for given user at given time
type Training struct {
	ID        string             `json:"id"`
	UserID    string             `json:"userId"`
	StartTime time.Time          `json:"startTime"`
	EndTime   time.Time          `json:"endTime,omitempty"`
	Exercises []TrainingExercise `json:"exercises"`
	Comment   string             `json:"comment"`
	CreatedAt time.Time          `json:"createdAt"`
}

// TrainingExercise keeps information about an exercise in the training
type TrainingExercise struct {
	ID         string        `json:"id"`
	ExerciseID string        `json:"exerciseId"`
	StartTime  time.Time     `json:"startTime"`
	EndTime    time.Time     `json:"endTime,omitempty"`
	Sets       []TrainingSet `json:"sets"`
	Comment    string        `json:"comment"`
	CreatedAt  time.Time     `json:"createdAt"`
}

// TrainingSet keeps information about a sets in the training
type TrainingSet struct {
	ID        string    `json:"id"`
	Time      time.Time `json:"time"`
	Reps      int       `json:"reps"`
	CreatedAt time.Time `json:"createdAt"`
}
