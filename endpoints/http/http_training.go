package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/unnamedxaer/gymm-api/entities"
	"github.com/unnamedxaer/gymm-api/repositories"
)

// StartTraining is a handler that trigger starting of a new training for logged in user.
func (app *App) StartTraining(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userID, ok := ctx.Value(contextKeyUserID).(string)
	if !ok {
		clearCookieJWTAuthToken(w)
		responseWithUnauthorized(w)
		return
	}

	tr, err := app.trainingUsecases.StartTraining(ctx, userID)
	if err != nil {
		logDebugError(app.l, req, err)
		var e *repositories.InvalidIDError
		if errors.As(err, &e) {
			responseWithError(w, http.StatusBadRequest, e)
			return
		}

		responseWithInternalError(w)
		return
	}

	responseWithJSON(w, http.StatusCreated, &tr)
}

// EndTraining is a handler that trigger starting of a new training for logged in user.
func (app *App) EndTraining(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userID, ok := ctx.Value(contextKeyUserID).(string)
	if !ok {
		clearCookieJWTAuthToken(w)
		responseWithUnauthorized(w)
		return
	}

	vars := mux.Vars(req)
	trainingID := vars["trainingID"]
	tr, err := app.trainingUsecases.GetTrainingByID(ctx, trainingID)
	if err != nil {
		logDebugError(app.l, req, err)
		var e *repositories.InvalidIDError
		if errors.As(err, &e) {
			responseWithError(w, http.StatusBadRequest, e)
			return
		}

		responseWithInternalError(w)
		return
	}

	if tr == nil || tr.UserID != userID {
		err = formatUnauthorizedError("training")
		logDebugError(app.l, req, err)
		responseWithError(w, http.StatusUnauthorized, err)
	}

	if !tr.EndTime.IsZero() {
		err := fmt.Errorf("training already completed")
		logDebugError(app.l, req, err)
		responseWithError(w, http.StatusConflict, err)
	}

	tr, err = app.trainingUsecases.EndTraining(ctx, trainingID)
	if err != nil {
		logDebugError(app.l, req, err)
		var e *repositories.InvalidIDError
		if errors.As(err, &e) {
			responseWithError(w, http.StatusBadRequest, e)
			return
		}

		responseWithInternalError(w)
		return
	}

	responseWithJSON(w, http.StatusOK, &tr)
}

// GetTrainingByID is a handler that returns user training for given id
func (app *App) GetTrainingByID(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userID, ok := ctx.Value(contextKeyUserID).(string)
	if !ok {
		clearCookieJWTAuthToken(w)
		responseWithUnauthorized(w)
		return
	}

	vars := mux.Vars(req)
	trainingID := vars["trainingID"]
	tr, err := app.trainingUsecases.GetTrainingByID(ctx, trainingID)
	if err != nil {
		logDebugError(app.l, req, err)
		var e *repositories.InvalidIDError
		if errors.As(err, &e) {
			responseWithError(w, http.StatusBadRequest, e)
			return
		}

		responseWithInternalError(w)
		return
	}

	if tr != nil && tr.UserID != userID {
		err = formatUnauthorizedError("training")
		logDebugError(app.l, req, err)
		responseWithError(w, http.StatusUnauthorized, err)
	}

	responseWithJSON(w, http.StatusOK, &tr)
}

// GetTraining is a handler that returns user trainings
func (app *App) GetUserTrainings(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userID, ok := ctx.Value(contextKeyUserID).(string)
	if !ok {
		clearCookieJWTAuthToken(w)
		responseWithUnauthorized(w)
		return
	}

	tr, err := app.trainingUsecases.GetUserTrainings(ctx, userID, false)
	if err != nil {
		logDebugError(app.l, req, err)
		var e *repositories.InvalidIDError
		if errors.As(err, &e) {
			responseWithError(w, http.StatusBadRequest, e)
			return
		}

		responseWithInternalError(w)
		return
	}

	responseWithJSON(w, http.StatusOK, &tr)
}

// StartTrainingExercise is a handler that adds new exercise to  the training
func (app *App) StartTrainingExercise(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userID, ok := ctx.Value(contextKeyUserID).(string)
	if !ok {
		clearCookieJWTAuthToken(w)
		responseWithUnauthorized(w)
		return
	}

	body := make(map[string]interface{}, 1)

	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		logDebugError(app.l, req, err)
		responseWithError(w, http.StatusBadRequest, err)
		return
	}

	exerciseID, ok := body["exerciseId"]
	if !ok {
		err := fmt.Errorf("missing %q property", "exerciseId")
		logDebugError(app.l, req, err)
		responseWithError(w, http.StatusBadRequest, err)
		return
	}

	exID, ok := exerciseID.(string)
	if !ok {
		err := fmt.Errorf(
			"incorrect type of %q property, expected string", "exerciseId")
		logDebugError(app.l, req, err)
		responseWithError(w, http.StatusBadRequest, err)
		return
	}

	vars := mux.Vars(req)
	trainingID := vars["trainingID"]
	tr, err := app.trainingUsecases.GetTrainingByID(ctx, trainingID)
	if err != nil {
		logDebugError(app.l, req, err)
		var e *repositories.InvalidIDError
		if errors.As(err, &e) {
			responseWithError(w, http.StatusBadRequest, e)
			return
		}

		responseWithInternalError(w)
		return
	}

	if tr == nil || tr.UserID != userID {
		err = formatUnauthorizedError("training")
		logDebugError(app.l, req, err)
		responseWithError(w, http.StatusUnauthorized, err)
		return
	}

	exercise, err := app.exerciseUsecases.GetExerciseByID(ctx, exID)
	if err != nil {
		logDebugError(app.l, req, err)
		var e *repositories.InvalidIDError
		if errors.As(err, &e) {
			responseWithError(w, http.StatusBadRequest, e)
			return
		}

		responseWithInternalError(w)
		return
	}

	if exercise == nil {
		err = fmt.Errorf("exercise with id %q does not exist", exID)
		logDebugError(app.l, req, err)
		responseWithError(w, http.StatusBadRequest, err)
		return
	}

	te := &entities.TrainingExercise{
		StartTime:  time.Now(),
		ExerciseID: exID,
	}

	te, err = app.trainingUsecases.StartExercise(ctx, tr.ID, te)
	if err != nil {
		logDebugError(app.l, req, err)
		var e *repositories.InvalidIDError
		if errors.As(err, &e) {
			responseWithError(w, http.StatusBadRequest, e)
			return
		}

		responseWithInternalError(w)
		return
	}

	responseWithJSON(w, http.StatusCreated, &te)
}

// EndTrainingExercise is a handler that stops training exercise
func (app *App) EndTrainingExercise(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userID, ok := ctx.Value(contextKeyUserID).(string)
	if !ok {
		clearCookieJWTAuthToken(w)
		responseWithUnauthorized(w)
		return
	}

	vars := mux.Vars(req)
	teID := vars["exerciseID"]

	te, err := app.trainingUsecases.EndExercise(ctx, userID, teID, time.Now())
	if err != nil {
		logDebugError(app.l, req, err)
		var e *repositories.InvalidIDError
		if errors.As(err, &e) {
			responseWithError(w, http.StatusBadRequest, e)
			return
		}

		responseWithInternalError(w)
		return
	}
	if te == nil {
		err = formatUnauthorizedError("training exercise")
		logDebugError(app.l, req, err)
		responseWithError(w, http.StatusUnauthorized, err)
		return
	}

	responseWithJSON(w, http.StatusOK, &te)
}

// AddTrainingSetExercise is a handler that adds new set to the training exercise
func (app *App) AddTrainingSetExercise(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userID, ok := ctx.Value(contextKeyUserID).(string)
	if !ok {
		clearCookieJWTAuthToken(w)
		responseWithUnauthorized(w)
		return
	}

	set := entities.TrainingSet{}

	err := json.NewDecoder(req.Body).Decode(&set)
	if err != nil {
		logDebugError(app.l, req, err)
		responseWithError(w, http.StatusBadRequest, err)
		return
	}

	vars := mux.Vars(req)
	teID := vars["exerciseID"]

	ts, err := app.trainingUsecases.AddSet(ctx, userID, teID, &set)
	if err != nil {
		logDebugError(app.l, req, err)
		var e *repositories.InvalidIDError
		if errors.As(err, &e) {
			responseWithError(w, http.StatusBadRequest, e)
			return
		}

		responseWithInternalError(w)
		return
	}

	responseWithJSON(w, http.StatusCreated, ts)
}
