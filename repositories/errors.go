package repositories

import (
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

// EmailAddressInUse is an error returned when user tries to register new account using email address that already exists in storage
type EmailAddressInUse struct {
	msg string
}

func (err EmailAddressInUse) Error() string {
	return err.msg
}

// NewErrorEmailAddressInUse returns a new error of type EmailAddressInUse
func NewErrorEmailAddressInUse() EmailAddressInUse {
	return EmailAddressInUse{
		msg: "email address already in use",
	}
}

// IsDuplicatedError checks whether given mongo error says that an insert violated unique contrain
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
