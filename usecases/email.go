package usecases

import (
	"bytes"
	"context"
	"net/smtp"
	"os"
	"text/template"

	"github.com/unnamedxaer/gymm-api/entities"
)

func sendResetPwdRequestEmail(ctx context.Context, user *entities.User) error {

	tmpl, err := template.ParseFiles("../templates/resetpwd.html")
	if err != nil {
		return err
	}

	b := bytes.Buffer{}
	err = tmpl.Execute(&b, user)
	if err != nil {
		return err
	}

	return sendEmail([]string{user.EmailAddress}, b.Bytes())
}

func sendEmail(recipients []string, data []byte) error {
	// Sender data.
	from := os.Getenv("APP_EMAIL_ADDRESS")
	password := os.Getenv("APP_EMAIL_PASSWORD")

	// Receiver email address.
	to := recipients

	// smtp server configuration.
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, data)
	if err != nil {
		return err
	}
	return nil
}
