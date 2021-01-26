package entities

import "time"

// User represents a person that uses the service
type User struct {
	ID           string
	FirstName    string
	LastName     string
	EmailAddress string
	CreatedAt    time.Time
}
