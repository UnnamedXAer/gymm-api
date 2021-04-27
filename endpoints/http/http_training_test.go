package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/unnamedxaer/gymm-api/mocks"
)

func TestTrainingHandlersUnauthorized(t *testing.T) {
	testCases := []struct {
		desc   string
		url    string
		method string
	}{
		{"get user trainings",
			"/trainings",
			http.MethodGet},

		{"start training",
			"/trainings/start",
			http.MethodPost},

		{"end training",
			"/trainings/" + mocks.ExampleTraining.ID + "/end",
			http.MethodPatch},

		{"get training by id",
			"/trainings/" + mocks.ExampleTraining.ID,
			http.MethodGet},

		{"start exercise",
			"/trainings/" + mocks.ExampleTraining.ID + "/exercises",
			http.MethodPost},

		{"end exercise",
			"/trainings/" + mocks.ExampleTraining.ID + "/exercises/end",
			http.MethodPatch},

		{"add set",
			"/trainings/" + mocks.ExampleTraining.ID + "/exercises/" + mocks.ExampleTraining.Exercises[0].ID + "/sets",
			http.MethodPatch},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			req, _ := http.NewRequest(tC.method, tC.url, nil)
			res := executeRequestWithoutJWT(req)
			checkResponseCode(t, http.StatusUnauthorized, res.Code)
		})
	}
}

func TestGetTrainings(t *testing.T) {

	req, _ := http.NewRequest(http.MethodGet, "/trainings/", nil)

	res := executeRequest(req)

	checkResponseCode(t, http.StatusOK, res.Code)

	if strings.Contains(res.Body.String(), "error") ||
		!strings.Contains(res.Body.String(), mocks.ExampleTraining.ID) {
		t.Errorf("want user trainings, got %s", res.Body.String())
	}
}

func TestStartTraining(t *testing.T) {

	req, _ := http.NewRequest(http.MethodPost, "/trainings/", nil)

	res := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, res.Code)

	if !strings.Contains(res.Body.String(), mocks.ExampleTraining.ID) {
		t.Errorf("want receive started training, got %s", res.Body.String())
	}
}

func TestEndTraining(t *testing.T) {

	req, _ := http.NewRequest(http.MethodPatch, "/trainings/"+mocks.ExampleTraining.ID+"/end", nil)

	res := executeRequest(req)

	checkResponseCode(t, http.StatusOK, res.Code)

	if !strings.Contains(res.Body.String(), mocks.ExampleTraining.ID) {
		t.Errorf("want receive ended training, got %s", res.Body.String())
	}
}

func TestStartExercise(t *testing.T) {

	payload := bytes.Buffer{}
	err := json.NewEncoder(&payload).Encode(mocks.ExampleExercise)
	if err != nil {
		t.Fatal(err)
	}

	req, _ := http.NewRequest(http.MethodPost, "/trainings/"+mocks.ExampleTraining.ID+"/exercises", &payload)

	res := executeRequest(req)

	checkResponseCode(t, http.StatusOK, res.Code)

	if !strings.Contains(res.Body.String(), mocks.ExampleExercise.ID) {
		t.Errorf("want receive started exercise, got %s", res.Body.String())
	}
}

func TestEndExercise(t *testing.T) {

	req, _ := http.NewRequest(http.MethodPatch,
		fmt.Sprintf("/trainings/%s/exercises/%s/end",
			mocks.ExampleTraining.ID, mocks.ExampleTraining.Exercises[0].ID),
		nil)

	res := executeRequest(req)

	checkResponseCode(t, http.StatusOK, res.Code)

	if !strings.Contains(res.Body.String(), mocks.ExampleTraining.Exercises[0].ID) {
		t.Errorf("want receive ended exercise, got %s", res.Body.String())
	}
}

func TestAddSet(t *testing.T) {

	payload := bytes.Buffer{}
	err := json.NewEncoder(&payload).Encode(mocks.ExampleTrainingSet)
	if err != nil {
		t.Fatal(err)
	}

	req, _ := http.NewRequest(http.MethodPost,
		fmt.Sprintf("/trainings/%s/exercises/%s/sets",
			mocks.ExampleTraining.ID, mocks.ExampleTraining.Exercises[0].ID),
		nil)

	res := executeRequest(req)

	checkResponseCode(t, http.StatusOK, res.Code)

	if !strings.Contains(res.Body.String(), mocks.ExampleTrainingSet.ID) {
		t.Errorf("want receive added set, got %s", res.Body.String())
	}
}
