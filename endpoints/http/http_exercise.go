package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/unnamedxaer/gymm-api/entities"
	"github.com/unnamedxaer/gymm-api/mocks"
	"github.com/unnamedxaer/gymm-api/repositories"
	"github.com/unnamedxaer/gymm-api/usecases"
	"github.com/unnamedxaer/gymm-api/validation"
)

func (app *App) CreateExercise(w http.ResponseWriter, req *http.Request) {
	var input usecases.ExerciseInput

	err := json.NewDecoder(req.Body).Decode(&input)
	if err != nil {
		app.l.Debug().Msgf("[POST / CreateExercise] -> body: %v, error: %v", input, err)

		syntaxErr, ok := err.(*json.SyntaxError)
		if ok {
			responseWithErrorMsg(w, http.StatusBadRequest, syntaxErr)
			return
		}

		invalidUnmarshalErr, ok := err.(*json.InvalidUnmarshalError)
		if ok {
			responseWithErrorMsg(w, http.StatusBadRequest, invalidUnmarshalErr)
			return
		}

		unmarshalTypeErr, ok := err.(*json.UnmarshalTypeError)
		if ok {
			responseWithErrorMsg(w, http.StatusBadRequest, unmarshalTypeErr)
			return
		}

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
			responseWithErrorJSON(w, http.StatusNotAcceptable, svErr.Format())
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

	vars := mux.Vars(req)

	id, ok := vars["id"]
	if !ok {
		err := errors.New("missign query parameter 'ID'")
		app.l.Debug().Msgf("update exercise: %v", err.Error())

		responseWithErrorMsg(w, http.StatusBadRequest, err)
		return
	}

	var input usecases.ExerciseInput
	err := json.NewDecoder(req.Body).Decode(&input)
	if err != nil {
		excludedTags := getExerciseExcludeTags(&input)
		errText := getErrOfMalformedInput(&input, excludedTags)

		err = errors.New(errText)
		app.l.Debug().Msgf("update exercise: %v", err.Error())
		responseWithErrorMsg(w, http.StatusBadRequest, err)
		return
	}
	defer req.Body.Close()
	app.l.Trace().Msgf("[Patch /exercise] -> body: %v", input)

	trimWhitespacesOnExerciseInput(&input)

	err = validateExerciseInput4Update(app.Validate, &input)
	if err != nil {
		app.l.Debug().Msgf("update exercise: %v", err.Error())

		if svErr, ok := err.(*validation.StructValidError); ok {
			responseWithErrorJSON(w, http.StatusNotAcceptable, svErr.Format())
			return
		}
		responseWithErrorMsg(w, http.StatusInternalServerError, err)
		return
	}

	curExercise, err := app.exerciseUsecases.GetExerciseByID(id)
	if err != nil {
		app.l.Debug().Msgf("update exercise: %v", err.Error())

		// @improvement: we could wrap original error with new one for
		// @improvement: logging purposes
		if !errors.Is(err, repositories.NewErrorNotFoundRecord()) {
			if errors.Is(err, repositories.NewErrorInvalidID(id)) {
				responseWithErrorJSON(w, http.StatusBadRequest, err)
				return
			}
			responseWithErrorJSON(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError))
			return
		}
	}

	if curExercise == nil || (curExercise.CreatedBy != mocks.UserID) { // @todo: authenticated user !!!
		err = fmt.Errorf("unauthorized: you do not have permissons to modify exercise with id %q", id)
		app.l.Debug().Msgf("update exercise: %v", err.Error())

		responseWithErrorMsg(w, http.StatusUnauthorized, err)
		return
	}

	exercise, err := app.exerciseUsecases.UpdateExercise(&entities.Exercise{
		ID:          id,
		Name:        input.Name,
		Description: input.Description,
		SetUnit:     input.SetUnit,
	})
	if err != nil {
		app.l.Debug().Msgf("update exercise: %v", err.Error())

		if repositories.IsDuplicatedError(err) {
			responseWithErrorMsg(w, http.StatusConflict, err) // @todo: create new error type
			return
		}

		responseWithErrorMsg(w, http.StatusInternalServerError, err)
		return
	}

	app.l.Trace().Msgf("[PATH /exercise] -> response: %v", exercise)
	responseWithJSON(w, http.StatusOK, exercise)
}
