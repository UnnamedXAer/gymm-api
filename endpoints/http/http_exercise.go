package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/unnamedxaer/gymm-api/entities"
	"github.com/unnamedxaer/gymm-api/repositories"
	"github.com/unnamedxaer/gymm-api/usecases"
	"github.com/unnamedxaer/gymm-api/validation"
)

func (app *App) CreateExercise(w http.ResponseWriter, req *http.Request) {
	var input usecases.ExerciseInput

	err := json.NewDecoder(req.Body).Decode(&input)
	if err != nil {
		app.l.Debug().Msgf("[POST / CreateExercise] -> body: %v, error: %v", input, err)
		excludedTags := getExerciseExcludeTags(&input)
		errText := getErrOfMalformedInput(&input, excludedTags)
		responseWithErrorMsg(w, http.StatusBadRequest, errors.New(errText))
		return
	}
	defer req.Body.Close()
	app.l.Debug().Msgf("[POST / CreateExercise] -> body: %v", input)

	trimWhitespacesOnExerciseInput(&input)

	err = validateExerciseInput(app.Validate, &input)
	if err != nil {
		if svErr, ok := err.(*validation.StructValidError); ok {
			responseWithErrorJSON(w, http.StatusNotAcceptable, svErr.ValidationErrors())
			return
		}
		responseWithErrorMsg(w, http.StatusInternalServerError, err)
		return
	}

	exercise, err := app.exerciseUsecases.CreateExercise(input.Name, input.Description, input.SetUnit, input.CreatedBy)
	if err != nil {
		if repositories.IsDuplicatedError(err) {
			responseWithErrorMsg(w, http.StatusConflict, err)
			return
		}

		responseWithErrorMsg(w, http.StatusInternalServerError, err)
		return
	}

	responseWithJSON(w, http.StatusCreated, exercise)
}

func (app *App) GetExeriseByID(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	id := vars["id"]
	app.l.Debug().Msg("[GET / GetExeriseByID] -> id: " + id)

	e, err := app.exerciseUsecases.GetExerciseByID(id)
	if err != nil {
		if errors.Is(err, repositories.NewErrorNotFoundRecord()) {
			responseWithJSON(w, http.StatusOK, nil)
			return
		}

		if errors.Is(err, repositories.NewErrorInvalidID(id)) {
			responseWithErrorMsg(w, http.StatusBadRequest, err)
			return
		}

		responseWithErrorMsg(w, http.StatusInternalServerError, err)
		return
	}

	responseWithJSON(w, http.StatusOK, e)
}

func (app *App) UpdateExercise(w http.ResponseWriter, req *http.Request) {

	responseWithErrorMsg(w, http.StatusNotImplemented, fmt.Errorf(http.StatusText(http.StatusNotImplemented)))
	return

	id := req.URL.Query().Get("id")

	var input usecases.ExerciseInput
	err := json.NewDecoder(req.Body).Decode(&input)
	if err != nil {
		app.l.Debug().Msgf("[POST / UpdateExercise] -> body: %v, error: %v", input, err)
		excludedTags := getExerciseExcludeTags(&input)
		errText := getErrOfMalformedInput(&input, excludedTags)
		responseWithErrorMsg(w, http.StatusBadRequest, errors.New(errText))
		return
	}
	defer req.Body.Close()
	app.l.Debug().Msgf("[POST / UpdateExercise] -> body: %v", input)

	trimWhitespacesOnExerciseInput(&input)

	err = validateExerciseInput(app.Validate, &input)
	// err = validateExerciseInput(app.Validate, &input)
	if err != nil {
		if svErr, ok := err.(*validation.StructValidError); ok {
			responseWithErrorJSON(w, http.StatusNotAcceptable, svErr.ValidationErrors())
			return
		}
		responseWithErrorMsg(w, http.StatusInternalServerError, err)
		return
	}

	exercise, err := app.exerciseUsecases.UpdateExercise(&entities.Exercise{
		ID:          id,
		Name:        input.Name,
		Description: input.Description,
		SetUnit:     input.SetUnit,
	})
	if err != nil {
		if repositories.IsDuplicatedError(err) {
			responseWithErrorMsg(w, http.StatusConflict, err)
			return
		}

		responseWithErrorMsg(w, http.StatusInternalServerError, err)
		return
	}

	responseWithJSON(w, http.StatusCreated, exercise)
}
