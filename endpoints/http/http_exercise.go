package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/unnamedxaer/gymm-api/repositories"
	"github.com/unnamedxaer/gymm-api/usecases"
	"github.com/unnamedxaer/gymm-api/validation"
)

func (app *App) CreateExercise(w http.ResponseWriter, req *http.Request) {
	var e usecases.ExerciseInput

	err := json.NewDecoder(req.Body).Decode(&e)
	if err != nil {
		app.l.Debug().Msgf("[POST / CreateExercise] -> body: %v, error: %v", e, err)
		excludedTags := make([]string, 2)
		tag, _ := validation.GetFieldJSONTag(&e, "CreatedBy")
		excludedTags[0] = tag
		tag, _ = validation.GetFieldJSONTag(&e, "CreatedAt")
		excludedTags[1] = tag
		errText := getErrOfMalformedInput(&e, excludedTags)
		responseWithErrorMsg(w, http.StatusBadRequest, errors.New(errText))
		return
	}
	defer req.Body.Close()
	app.l.Debug().Msgf("[POST / CreateExercise] -> body: %v", e)

	trimWhitespacesOnExerciseInput(&e)

	err = validateExerciseInput(app.Validate, &e)
	if err != nil {
		if svErr, ok := err.(*validation.StructValidError); ok {
			responseWithErrorJSON(w, http.StatusNotAcceptable, svErr.ValidationErrors())
			return
		}
		responseWithErrorMsg(w, http.StatusInternalServerError, err)
		return
	}

	exercise, err := app.exerciseUsecases.CreateExercise(e.Name, e.Description, e.SetUnit, e.CreatedBy)
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
