package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/unnamedxaer/gymm-api/entities"
	"github.com/unnamedxaer/gymm-api/mocks"
	"github.com/unnamedxaer/gymm-api/usecases"
)

func TestLogin(t *testing.T) {
	u := usecases.UserInput{
		EmailAddress: mocks.ExampleUser.EmailAddress,
		Password:     string(mocks.Password),
	}

	uJSON, err := json.Marshal(&u)
	if err != nil {
		t.Fatalf("could not marshal login payload: %v", err)
	}

	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(uJSON))
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	res := response.Result()
	cookies := res.Cookies()

	var tokenCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == cookieJwtTokenName {
			tokenCookie = c
		}
	}

	if tokenCookie == nil {
		t.Errorf("want token cookie, got nil")
	} else {

		if !tokenCookie.HttpOnly {
			t.Errorf("want token cookie to be HttpOnly")
		}

		if !tokenCookie.Expires.IsZero() {
			t.Errorf("want token cookie expire time to be zero, got %v", tokenCookie.Expires)
		}
	}

	var got struct {
		Error string
		User  *entities.User
	}
	err = json.NewDecoder(response.Body).Decode(&got)
	if err != nil {
		t.Fatalf("could not decode login response: %v", err)
	}

	if got.User == nil || got.Error != "" {
		t.Errorf("want authenticated user and no error, got %#v", got)
	}

}

func TestRegister(t *testing.T) {
	u := correctUser

	uJSON, _ := json.Marshal(u)

	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(uJSON))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)
}

func TestRegisterMalformedData(t *testing.T) {
	uStr := `{"userName: "Al", "emailAddress" : "al@mymeil.com", "password":"pwd123"}`
	uJSON := []byte(uStr)

	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(uJSON))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusUnprocessableEntity, response.Code)
}

func TestRegisterValidationFail(t *testing.T) {
	u := wrongUser

	uJSON, _ := json.Marshal(u)

	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(uJSON))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotAcceptable, response.Code)
}
