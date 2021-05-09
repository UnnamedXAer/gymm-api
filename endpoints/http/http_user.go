package http

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unnamedxaer/gymm-api/usecases"
)

func (app *App) GetUserById(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	id := vars["id"]
	logDebug(app.l, req, id)
	if id == "" {
		responseWithError(w, http.StatusUnprocessableEntity, errors.New("missign 'ID'"))
		return
	}

	ctx := req.Context()

	u, err := app.userUsecases.GetUserByID(ctx, id)
	if err != nil {
		var e *usecases.InvalidIDError
		if errors.As(err, &e) {
			responseWithError(w, http.StatusBadRequest, err)
			return
		}

		responseWithError(w, http.StatusInternalServerError, err)
		return
	}

	responseWithJSON(w, http.StatusOK, u)
}
