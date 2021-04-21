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

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(uJSON))
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var got struct {
		Error string
		User  *entities.User
	}
	err = json.NewDecoder(response.Body).Decode(&got)
	if err != nil {
		t.Fatalf("could not decode login response: %v", err)
	}

	if got.User == nil || got.Error != "" {
		t.Fatalf("want authenticated user and no error, got %#v", got)
	}

}
