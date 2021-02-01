package http

import (
	"strings"
	"testing"

	"github.com/unnamedxaer/gymm-api/usecases"
)

func TestValidateUserInputWrongInput(t *testing.T) {
	err := validateUserInput(validate, &wrongUser)

	if err == nil {
		t.Fatalf("Expected to get error, got %s \nuser: %+v", err, wrongUser)
	}

	if strings.Contains(err.Error(), "'emailAddress'") == false {
		t.Fatalf("Expected to get error about wrong field 'emailAddress', got '%s'", err)
	}
	if strings.Contains(err.Error(), "'userName'") == false {
		t.Fatalf("Expected to get error about wrong field 'userName', got '%s'", err)
	}
	if strings.Contains(err.Error(), "'password'") == false {
		t.Fatalf("Expected to get error about wrong field 'password', got '%s'", err)
	}
}

func TestValidateUserInputWrongEmail(t *testing.T) {
	err := validateUserInput(validate, &usecases.UserInput{
		EmailAddress: wrongUser.EmailAddress,
		Username:     correctUser.Username,
		Password:     correctUser.Password,
	})

	if err == nil {
		t.Fatalf("Expected to get error, got %s \nuser: %+v", err, wrongUser)
	}
	if strings.Contains(err.Error(), "'emailAddress'") == false {
		t.Fatalf("Expected to get error about wrong field 'emailAddress', got '%s'", err)
	}
	if strings.Contains(err.Error(), "'userName'") || strings.Contains(err.Error(), "'password'") {
		t.Fatalf("Expected to got error about 'emailAddress', got %s", err.Error())
	}
}
func TestValidateUserInputWrongUsername(t *testing.T) {
	err := validateUserInput(validate, &usecases.UserInput{
		Username:     wrongUser.Username,
		EmailAddress: correctUser.EmailAddress,
		Password:     correctUser.Password,
	})

	if err == nil {
		t.Fatalf("Expected to get error, got %s \nuser: %+v", err, wrongUser)
	}
	if strings.Contains(err.Error(), "'userName'") == false {
		t.Fatalf("Expected to get error about wrong field 'userName', got '%s'", err)
	}
	if strings.Contains(err.Error(), "'emailAddress'") || strings.Contains(err.Error(), "'password'") {
		t.Fatalf("Expected to got error about 'userName', got %s", err.Error())
	}
}

func TestValidateUserInputWrongPassword(t *testing.T) {
	err := validateUserInput(validate, &usecases.UserInput{
		Password:     wrongUser.Password,
		Username:     correctUser.Username,
		EmailAddress: correctUser.EmailAddress,
	})

	if err == nil {
		t.Fatalf("Expected to get error, got %s \nuser: %+v", err, wrongUser)
	}
	if strings.Contains(err.Error(), "'password'") == false {
		t.Fatalf("Expected to get error about wrong field 'password', got '%s'", err)
	}
	if strings.Contains(err.Error(), "'userName'") || strings.Contains(err.Error(), "'emailAddress'") {
		t.Fatalf("Expected to got error about 'password', got %s", err.Error())
	}
}

func TestValidateUserInputCorrectInput(t *testing.T) {
	err := validateUserInput(validate, &correctUser)

	if err != nil {
		t.Fatalf("Expected to pass validation, got error: '%s' \nuser: %+v", err, correctUser)
	}
}
