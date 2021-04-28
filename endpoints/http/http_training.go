package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/unnamedxaer/gymm-api/entities"
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

	tr, err := app.trainingUsecases.StartTraining(userID)
	if err != nil {
		logDebugError(app.l, req, err)
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
	trainingID := vars["trainingId"]
	tr, err := app.trainingUsecases.GetTrainingByID(trainingID)
	if err != nil {
		logDebugError(app.l, req, err)
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

	tr, err = app.trainingUsecases.EndTraining(trainingID)
	if err != nil {
		logDebugError(app.l, req, err)
		responseWithInternalError(w)
		return
	}

	responseWithJSON(w, http.StatusCreated, &tr)
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
	trainingID := vars["trainingId"]
	tr, err := app.trainingUsecases.GetTrainingByID(trainingID)
	if err != nil {
		logDebugError(app.l, req, err)
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

	tr, err := app.trainingUsecases.GetUserTrainings(userID, false)
	if err != nil {
		logDebugError(app.l, req, err)
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

	vars := mux.Vars(req)

	exerciseID, ok := vars["exerciseId"]
	if !ok {
		err := fmt.Errorf("missing %q property", "exerciseId")
		logDebugError(app.l, req, err)
		responseWithError(w, http.StatusBadRequest, err)
		return
	}

	trainingID := vars["trainingId"]
	tr, err := app.trainingUsecases.GetTrainingByID(trainingID)
	if err != nil {
		logDebugError(app.l, req, err)
		responseWithInternalError(w)
		return
	}

	if tr == nil || tr.UserID != userID {
		err = formatUnauthorizedError("training")
		logDebugError(app.l, req, err)
		responseWithError(w, http.StatusUnauthorized, err)
	}

	// @todo: verify if exercise with such id exists.
	responseWithErrorTxt(w, http.StatusNotImplemented, http.StatusText(http.StatusNotImplemented))
	return
	te := &entities.TrainingExercise{
		StartTime:  time.Now(),
		ExerciseID: exerciseID,
	}

	te, err = app.trainingUsecases.AddExercise(tr.ID, te)
	if err != nil {
		logDebugError(app.l, req, err)
		responseWithInternalError(w)
		return
	}

	responseWithJSON(w, http.StatusCreated, &tr)
}

// EndTrainingExercise is a handler that stops training exercise
func (app *App) EndTrainingExercise(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	w.Write([]byte(fmt.Sprintf("[%s] %s - (EndTrainingExercise) not implemented yet \n \n %v", req.Method, req.RequestURI, vars)))
}

// AddTrainingSetExercise is a handler that adds new set to the training exercise
func (app *App) AddTrainingSetExercise(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	w.Write([]byte(fmt.Sprintf("[%s] %s - (AddTrainingSetExercise) not implemented yet \n \n %v", req.Method, req.RequestURI, vars)))
}
