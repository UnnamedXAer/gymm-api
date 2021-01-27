package repositories

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
