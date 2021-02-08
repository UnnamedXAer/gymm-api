package repositories

import (
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

// EmailAddressInUseError is an error returned when user tries to register new account using email address that already exists in storage
type EmailAddressInUseError struct {
	msg string
}

func (err EmailAddressInUseError) Error() string {
	return err.msg
}

// NewErrorEmailAddressInUse returns a new error of type EmailAddressInUse
func NewErrorEmailAddressInUse() EmailAddressInUseError {
	return EmailAddressInUseError{
		msg: "email address already in use",
	}
}

// NotFoundRecordError is an error returned when single row result query did not found matching data
type NotFoundRecordError struct {
	msg string
}

func (err NotFoundRecordError) Error() string {
	return err.msg
}

// NewErrorNotFoundRecord returns a new error of type NotFoundRecord
func NewErrorNotFoundRecord() NotFoundRecordError {
	return NotFoundRecordError{
		msg: "record not found",
	}
}

// InvalidIDError is an error returned when given ID is not valid
type InvalidIDError struct {
	msg string
}

func (err InvalidIDError) Error() string {
	return err.msg
}

// NewErrorInvalidID returns a new error of type InvalidID
func NewErrorInvalidID(id string) InvalidIDError {
	return InvalidIDError{
		msg: "invalid ID: " + id,
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
