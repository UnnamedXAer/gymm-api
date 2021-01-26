package entities

import "time"

// Training keeps an informations about set of executed exercises for given user at given time
type Training struct {
	ID        string
	UserID    string
	StartTime time.Time
	EndTime   time.Time
	Exercises []interface{}
	Comment   string
}
