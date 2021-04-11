package entities

import "time"

type SetUnit int8

const (
	Weight SetUnit = iota + 1
	Time
)

type Exercise struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	SetUnit     SetUnit   `json:"setUnit"`
	CreatedAt   time.Time `json:"createdAt"`
	CreatedBy   string    `json:"createdBy"`
}
