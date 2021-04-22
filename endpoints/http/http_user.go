package http

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unnamedxaer/gymm-api/repositories"
)

func (app *App) GetUserById(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	id := vars["id"]
	app.l.Debug().Msg("[GET / GetUserById] -> id: " + id)
	if id == "" {
		responseWithError(w, http.StatusUnprocessableEntity, errors.New("missign 'ID'"))
		return
	}

	u, err := app.userUsecases.GetUserByID(id)
	if err != nil {
		var e *repositories.InvalidIDError
		if errors.As(err, &e) {
			responseWithError(w, http.StatusBadRequest, err)
			return
		}

		responseWithError(w, http.StatusInternalServerError, err)
		return
	}

	responseWithJSON(w, http.StatusOK, u)
}
