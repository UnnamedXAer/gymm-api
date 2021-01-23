package models

import (
	"time"
)

type User struct {
	ID           string `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname    string `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname     string `json:"lastname,omitempty" bson:"lastname,omitempty"`
	EmailAddress string `json:"emailaddress,omitempty" bson:"emailaddress,omitempty"`
	// Password     string    `json:"password,omitempty" bson:"password,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}
