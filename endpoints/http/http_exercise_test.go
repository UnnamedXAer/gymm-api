package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/unnamedxaer/gymm-api/entities"
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
			desc: "exercise with correct data",
			input: usecases.ExerciseInput{
				Name:        mocks.ExampleExercise.Name,
				Description: mocks.ExampleExercise.Description,
				SetUnit:     mocks.ExampleExercise.SetUnit,
			},
			want: http.StatusCreated,
		},
		{
			desc: "exercise with set unit as time",
			input: usecases.ExerciseInput{
				Name:        mocks.ExampleExercise.Name,
				Description: mocks.ExampleExercise.Description,
				SetUnit:     entities.Time,
			},
			want: http.StatusCreated,
		},
		{
			desc: "exercise witout Description",
			input: usecases.ExerciseInput{
				Name:    mocks.ExampleExercise.Name,
				SetUnit: mocks.ExampleExercise.SetUnit,
			},
			want: http.StatusNotAcceptable,
		},
		{
			desc: "exercise without Name",
			input: usecases.ExerciseInput{
				Description: mocks.ExampleExercise.Description,
				SetUnit:     mocks.ExampleExercise.SetUnit,
			},
			want: http.StatusNotAcceptable,
		},
		{
			desc: "exercise without SetUnit",
			input: usecases.ExerciseInput{
				Name:        mocks.ExampleExercise.Name,
				Description: mocks.ExampleExercise.Description,
			},
			want: http.StatusNotAcceptable,
		},
		{
			desc: "exercise with wrong SetUnit",
			input: usecases.ExerciseInput{
				Name:        mocks.ExampleExercise.Name,
				Description: mocks.ExampleExercise.Description,
				SetUnit:     123,
			},
			want: http.StatusNotAcceptable,
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

func TestCreateExerciseUnauthorized(t *testing.T) {
	payload := []byte(`{"name":"DL"}`)

	req, _ := http.NewRequest(http.MethodPatch, "/exercises/"+mocks.ExampleExercise.ID, bytes.NewBuffer(payload))

	res := executeRequestWithoutJWT(req)

	checkResponseCode(t, http.StatusUnauthorized, res.Code)
}

func TestCreateExerciseMalformedData(t *testing.T) {
	payload := []byte(`{"name:"Deadlift","description":"The deadlift is a ...","setUnit":1}`)

	req, _ := http.NewRequest(http.MethodPost, "/exercises", bytes.NewBuffer(payload))

	res := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, res.Code)
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

func TestGetExerciseUnauthorized(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/exercises/"+mocks.ExampleExercise.ID, nil)

	res := executeRequestWithoutJWT(req)

	checkResponseCode(t, http.StatusUnauthorized, res.Code)
}

func TestUpdateExercise(t *testing.T) {
	payload := []byte(`{"name":"DL"}`)

	req, _ := http.NewRequest(http.MethodPatch, "/exercises/"+mocks.ExampleExercise.ID, bytes.NewBuffer(payload))

	res := executeRequest(req)

	checkResponseCode(t, http.StatusOK, res.Code)
}

func TestUpdateExerciseUnauthorized(t *testing.T) {
	payload := []byte(`{"name":"DL"}`)

	req, _ := http.NewRequest(http.MethodPatch, "/exercises/"+mocks.ExampleExercise.ID, bytes.NewBuffer(payload))

	res := executeRequestWithoutJWT(req)

	checkResponseCode(t, http.StatusUnauthorized, res.Code)
}

func TestUpdateExerciseDifferentUserUnauthorized(t *testing.T) {
	payload := []byte(`{"name":"DL"}`)

	req, _ := http.NewRequest(http.MethodPatch, "/exercises/"+mocks.ExampleExercise.ID, bytes.NewBuffer(payload))

	mocks.ExampleExercise.CreatedBy += "1"
	res := executeRequest(req)
	mocks.ExampleExercise.CreatedBy = mocks.UserID[:len(mocks.UserID)-1]
	checkResponseCode(t, http.StatusUnauthorized, res.Code)
}

func TestUpdateExerciseMalformedData(t *testing.T) {
	payload := []byte(`{"name:"DL"}`)

	req, _ := http.NewRequest(http.MethodPatch, "/exercises/"+mocks.ExampleExercise.ID, bytes.NewBuffer(payload))

	res := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, res.Code)
}

func TestUpdateExerciseIncorrectID(t *testing.T) {
	payload := []byte(`{"name":"DL"}`)

	incorrectID := "12435678901234567890123z"

	req, err := http.NewRequest(http.MethodPatch, "/exercises/"+incorrectID, bytes.NewBuffer(payload))

	if err != nil {
		t.Fatal(err)
	}

	res := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, res.Code)

	// @todo: correct after response corrected uncomment it.
	// wantResponse := repositories.NewErrorInvalidID(incorrectID).Error()
	// if !strings.Contains(res.Body.String(), wantResponse) {
	// 	t.Fatalf("want to get response like %q, got %q", wantResponse, res.Body.String())
	// }
}

func TestGetExercisesByName(t *testing.T) {

	testCases := []struct {
		desc  string
		input string
		// want is a len of the unmarshalled response body
		want int
	}{
		{
			desc:  "existing exercise",
			input: strings.ToLower(mocks.ExampleExercise.Name),
			want:  1,
		},
		{
			desc:  "not existing exercise",
			input: "notfound - Exercise Name",
			want:  0,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/exercises?n="+tC.input, nil)
			if err != nil {
				t.Error(err)
				return
			}
			res := executeRequest(req)
			checkResponseCode(t, http.StatusOK, res.Code)

			exercises := []entities.Exercise{}
			err = json.Unmarshal(res.Body.Bytes(), &exercises)
			if err != nil {
				t.Error(err)
				return
			}
			if len(exercises) != tC.want {
				t.Errorf("want unmarshalled results len eq %d, got %d, for %q", tC.want, len(exercises), tC.input)
			}
		})
	}
}
func TestGetExercisesByNameMissingParam(t *testing.T) {
	want := `"error":"missing name`

	req, err := http.NewRequest(http.MethodGet, "/exercises", nil)
	if err != nil {
		t.Error(err)
		return
	}
	res := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, res.Code)

	got := res.Body.String()
	if !strings.Contains(got, want) {
		t.Errorf("want response like 'missing name...', got %q for req with no 'n' parameter", got)
	}

	req, err = http.NewRequest(http.MethodGet, "/exercises?n=", nil)
	if err != nil {
		t.Error(err)
		return
	}
	res = executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, res.Code)

	got = res.Body.String()
	if !strings.Contains(got, want) {
		t.Errorf("want response like 'missing name...', got %q for empty 'n' parameter", got)
	}
}

func TestGetExercisesByNameEmptyParam(t *testing.T) {
	want := `"error":"missing name`
	req, err := http.NewRequest(http.MethodGet, "/exercises?n=", nil)
	if err != nil {
		t.Error(err)
		return
	}
	res := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, res.Code)

	got := res.Body.String()
	if !strings.Contains(got, want) {
		t.Errorf("want response like 'missing name...', got %q for empty 'n' parameter", got)
	}
}
