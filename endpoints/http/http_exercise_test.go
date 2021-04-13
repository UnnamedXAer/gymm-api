package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/unnamedxaer/gymm-api/mocks"
	"github.com/unnamedxaer/gymm-api/usecases"
)

func TestCreateExercise(t *testing.T) {
	testCases := []struct {
		desc  string
		input usecases.ExerciseInput
		want  int
	}{
		{
			desc: "exercise based on example exerice data",
			input: usecases.ExerciseInput{
				Name:        mocks.ExampleExercise.Name,
				Description: mocks.ExampleExercise.Description,
				SetUnit:     mocks.ExampleExercise.SetUnit,
				CreatedBy:   mocks.UserID,
			},
			want: http.StatusCreated,
		},
		{
			desc: "exercise witout Description",
			input: usecases.ExerciseInput{
				Name:      mocks.ExampleExercise.Name,
				SetUnit:   mocks.ExampleExercise.SetUnit,
				CreatedBy: mocks.UserID,
			},
			want: http.StatusCreated,
		},
		{
			desc: "exercise without Name",
			input: usecases.ExerciseInput{
				Description: mocks.ExampleExercise.Description,
				SetUnit:     mocks.ExampleExercise.SetUnit,
				CreatedBy:   mocks.UserID,
			},
			want: http.StatusUnprocessableEntity,
		},
		{
			desc: "exercise without SetUnit",
			input: usecases.ExerciseInput{
				Name:        mocks.ExampleExercise.Name,
				Description: mocks.ExampleExercise.Description,
				CreatedBy:   mocks.UserID,
			},
			want: http.StatusUnprocessableEntity,
		},
		{
			desc: "exercise with wring SetUnit",
			input: usecases.ExerciseInput{
				Name:        mocks.ExampleExercise.Name,
				Description: mocks.ExampleExercise.Description,
				SetUnit:     123,
				CreatedBy:   mocks.UserID,
			},
			want: http.StatusUnprocessableEntity,
		},
		{
			desc: "exercise without UserID",
			input: usecases.ExerciseInput{
				Name:        mocks.ExampleExercise.Name,
				Description: mocks.ExampleExercise.Description,
				SetUnit:     mocks.ExampleExercise.SetUnit,
			},
			want: http.StatusUnauthorized,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			payload, err := json.Marshal(tC.input)
			if err != nil {
				t.Error(err)
				return
			}
			req, err := http.NewRequest(http.MethodPost, "/exercises", bytes.NewBuffer(payload))
			if err != nil {
				t.Error(err)
				return
			}

			res := executeRequest(req)
			checkResponseCode(t, tC.want, res.Code)
		})
	}
}

func TestCreateExerciseMalformedData(t *testing.T) {
	payload := []byte(`{"name:"Deadlift","description":"The deadlift is a weight training exercise in which a loaded barbell or bar is lifted off the ground to the level of the hips, torso perpendicular to the floor, before being placed back on the ground. It is one of the three powerlifting exercises, along with the squat and bench press.","setUnit":1}`)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusUnprocessableEntity, response.Code)
}

func TestGetExerciseByID(t *testing.T) {

	testCases := []struct {
		desc  string
		input string
		// want is a len of the unmarshalled response body
		want int
	}{
		{
			desc:  "existing exercise",
			input: mocks.ExampleExercise.ID,
			want:  6,
		},
		{
			desc:  "not existing exercise",
			input: "606ea1de1c4e78b2da793211",
			want:  0,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/exercises/"+tC.input, nil)
			if err != nil {
				t.Error(err)
				return
			}
			req.Header.Set("Content-Type", "application/json")
			res := executeRequest(req)
			checkResponseCode(t, http.StatusOK, res.Code)

			// mb is enough to call res.Body.String() = "xyz"
			ex := make(map[string]interface{})
			err = json.Unmarshal(res.Body.Bytes(), &ex)
			if err != nil {
				t.Error(err)
				return
			}
			if len(ex) != tC.want {
				t.Errorf("want unmarshalled res len eq %d, got %d, for %q", tC.want, len(ex), res.Body.String())
			}
		})
	}
}
