package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

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

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(uJSON))
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

		if time.Until(tokenCookie.Expires) < 0 {
			t.Errorf("want token expire time to be in the future")
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
