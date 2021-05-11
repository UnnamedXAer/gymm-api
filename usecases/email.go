package usecases

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/unnamedxaer/gymm-api/entities"
)

type Mailer interface {
	Send(recipients []string, subject, data []byte)
	Close()
}

func generatePwdResetEmailContent(
	user *entities.User,
	pwdResetReq *entities.ResetPwdReq) ([]byte, error) {

	tmpl, err := template.ParseFiles("../templates/templatefiles/resetpwd.html")
	if err != nil {
		return nil, err
	}

	clientURL := os.Getenv("CLIENT_URL")
	appName := "The Gymm Api"

	url := fmt.Sprintf("%s/password/reset/%s", clientURL, pwdResetReq.ID)

	data := map[string]interface{}{
		"User":    user,
		"AppName": appName,
		"URL":     url,
	}

	b := bytes.Buffer{}
	err = tmpl.Execute(&b, &data)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
