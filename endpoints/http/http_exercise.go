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
		logDebugError(app.l, req, err)

		ok, err := formatParseErrors(err)
		if ok {
			responseWithErrorMsg(w, http.StatusBadRequest, err)
			return
		}

		errText := getErrOfMalformedInput(&input, exerciseExcludedFields)

		responseWithErrorMsgTxt(w, http.StatusBadRequest, errText)
		return
	}
	defer req.Body.Close()
	app.l.Debug().Msgf("[POST / CreateExercise] -> body: %v", input)

	trimWhitespacesOnExerciseInput(&input)

	err = validateExerciseInput(app.Validate, &input)
	if err != nil {
		logDebugError(app.l, req, err)
		if svErr, ok := err.(*validation.StructValidError); ok {
			responseWithJSON(w, http.StatusNotAcceptable, svErr.Format())
			return
		}
		responseWithInternalError(w)
		return
	}

	exercise, err := app.exerciseUsecases.CreateExercise(input.Name, input.Description, input.SetUnit, mocks.UserID) // @todo: logged user
	if err != nil {
		logDebugError(app.l, req, err)
		if repositories.IsDuplicatedError(err) {
			responseWithErrorMsg(w, http.StatusConflict,
				fmt.Errorf("exercise with name: %q and set unit: %d already exists",
					input.Name, input.SetUnit))
			return
		}

		responseWithInternalError(w)
		return
	}

	responseWithJSON(w, http.StatusCreated, exercise)
}

func (app *App) GetExeriseByID(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	id := vars["id"]
	app.l.Debug().Msg("[GET / GetExeriseByID] -> id: " + id)

	exercise, err := app.exerciseUsecases.GetExerciseByID(id)
	if err != nil {
		logDebugError(app.l, req, err)
		var e *repositories.InvalidIDError
		if errors.As(err, &e) {
			responseWithErrorMsg(w, http.StatusBadRequest, e)
			return
		}

		responseWithInternalError(w)
		return
	}

	responseWithJSON(w, http.StatusOK, exercise)
}

func (app *App) UpdateExercise(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)

	id, ok := vars["id"]
	if !ok {
		err := errors.New("missign query parameter 'ID'")
		logDebugError(app.l, req, err)

		responseWithErrorMsg(w, http.StatusBadRequest, err)
		return
	}

	var input usecases.ExerciseInput
	err := json.NewDecoder(req.Body).Decode(&input)
	if err != nil {
		errText := getErrOfMalformedInput(&input, exerciseExcludedFields)

		err = errors.New(errText)
		logDebugError(app.l, req, err)
		responseWithErrorMsg(w, http.StatusBadRequest, err)
		return
	}
	defer req.Body.Close()
	app.l.Trace().Msgf("[Patch /exercise] -> body: %v", input)

	trimWhitespacesOnExerciseInput(&input)

	err = validateExerciseInput4Update(app.Validate, &input)
	if err != nil {
		logDebugError(app.l, req, err)

		if svErr, ok := err.(*validation.StructValidError); ok {
			responseWithJSON(w, http.StatusNotAcceptable, svErr.Format())
			return
		}
		responseWithInternalError(w)
		return
	}

	curExercise, err := app.exerciseUsecases.GetExerciseByID(id)
	if err != nil {
		logDebugError(app.l, req, err)
		var e *repositories.InvalidIDError
		if errors.As(err, &e) {
			responseWithErrorMsg(w, http.StatusBadRequest, e)
			return
		}
		responseWithInternalError(w)
		return
	}

	if curExercise == nil || (curExercise.CreatedBy != mocks.UserID) { // @todo: authenticated user !!!
		err = fmt.Errorf("unauthorized: you do not have permissons to modify exercise (%s)", id)
		logDebugError(app.l, req, err)

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
		logDebugError(app.l, req, err)

		if repositories.IsDuplicatedError(err) {
			responseWithErrorMsg(w, http.StatusConflict, err) // @todo: create new error type
			return
		}

		responseWithInternalError(w)
		return
	}

	app.l.Trace().Msgf("[PATCH /exercise] -> response: %v", exercise)
	responseWithJSON(w, http.StatusOK, exercise)
}
