package usecases

import (
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

// EmailAddressInUseError is an error returned when user tries to register new account or update his email using email address that already exists in storage
type EmailAddressInUseError struct {
}

func (err EmailAddressInUseError) Error() string {
	return "email address already in use"
}

// NewErrorEmailAddressInUse returns a new error of type EmailAddressInUse
func NewErrorEmailAddressInUse() *EmailAddressInUseError {
	return &EmailAddressInUseError{}
}

// InvalidIDError is an error returned when given ID is not valid
type InvalidIDError struct {
	ID       string
	dataName string
}

func (err InvalidIDError) Error() string {
	return "invalid " + err.dataName + " ID: " + err.ID
}

// NewErrorInvalidID returns a new error of type InvalidID
func NewErrorInvalidID(id string, dataName string) *InvalidIDError {
	return &InvalidIDError{
		ID:       id,
		dataName: dataName,
	}
}

// RecordNotExistsError is an error returned when expected record does not exist in storage
type RecordNotExistsError struct {
	dataName string
}

func (err RecordNotExistsError) Error() string {
	return err.dataName + " does not exist"
}

// NewErrorRecordNotExists returns a new error of type *RecordNotExistsError
func NewErrorRecordNotExists(dataName string) *RecordNotExistsError {
	return &RecordNotExistsError{
		dataName: dataName,
	}
}

// IsDuplicatedError checks whether given mongo error says that an insert violated unique constrain
func IsDuplicatedError(err error) bool {
	var e mongo.WriteException

	if errors.As(err, &e) {
		for _, we := range e.WriteErrors {
			if we.Code == 11000 {
				return true
			}
		}
	}
	return false
}
