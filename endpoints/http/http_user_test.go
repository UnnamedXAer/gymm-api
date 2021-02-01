package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestCreateUser(t *testing.T) {
	u := correctUser

	uJSON, _ := json.Marshal(u)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(uJSON))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)
}

func TestCreateUserMalformedData(t *testing.T) {
	uStr := `{"userName: "Al", "emailAddress" : "al@mymeil.com", "password":"pwd123"}`
	uJSON := []byte(uStr)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(uJSON))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusUnprocessableEntity, response.Code)
}

func TestCreateUserValidationFail(t *testing.T) {
	u := wrongUser

	uJSON, _ := json.Marshal(u)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(uJSON))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotAcceptable, response.Code)
}

func TestGetUserByID(t *testing.T) {
	id := "1sadf3245df3245"

	req, _ := http.NewRequest(http.MethodGet, ("/users/" + id), nil)
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	//check returned obj
}

func TestGetUserByIDBotFound(t *testing.T) {
	id := "1sadf3245df3245"

	req, _ := http.NewRequest(http.MethodGet, ("/users/" + id + "not_found"), nil)
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	//check returned obj == nil
}
