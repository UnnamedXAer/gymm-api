package entities

import "time"

// User represents a person that uses the service
type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"userName"`
	EmailAddress string    `json:"emailAddress"`
	CreatedAt    time.Time `json:"createdAt"`
}
