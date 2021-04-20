package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unnamedxaer/gymm-api/repositories"
	"github.com/unnamedxaer/gymm-api/usecases"
	"github.com/unnamedxaer/gymm-api/validation"
)

func (app *App) CreateUser(w http.ResponseWriter, req *http.Request) {
	var u usecases.UserInput
	f, ff := validation.GetFieldJSONTag(&u, "Username")
	fmt.Println(f, ff)
	err := json.NewDecoder(req.Body).Decode(&u)
	app.l.Debug().Msg("[POST / CreateUser] -> body: " + fmt.Sprintf("%v", u))
	if err != nil {
		resErrText := getErrOfMalformedInput(&u, []string{"ID", "CreatedAt"})
		responseWithErrorMsg(w, http.StatusUnprocessableEntity, errors.New(resErrText))
		return
	}
	defer req.Body.Close()

	trimWhitespacesOnUserInput(&u)

	err = validateUserInput(app.Validate, &u)
	if err != nil {
		if svErr, ok := err.(*validation.StructValidError); ok {
			responseWithJSON(w, http.StatusNotAcceptable, svErr.ValidationErrors())
			return
		}
		responseWithErrorMsg(w, http.StatusInternalServerError, err)
		return
	}

	user, err := app.userUsecases.CreateUser(&u)
	if err != nil {
		if errors.Is(err, repositories.NewErrorEmailAddressInUse()) {
			responseWithErrorMsg(w, http.StatusConflict, err)
			return
		}

		responseWithErrorMsg(w, http.StatusInternalServerError, err)
		return
	}

	responseWithJSON(w, http.StatusCreated, user)
}

func (app *App) GetUserById(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	id := vars["id"]
	app.l.Debug().Msg("[GET / GetUserById] -> id: " + id)
	if id == "" {
		responseWithErrorMsg(w, http.StatusUnprocessableEntity, errors.New("missign 'ID'"))
		return
	}

	u, err := app.userUsecases.GetUserByID(id)
	if err != nil {
		var e *repositories.InvalidIDError
		if errors.As(err, &e) {
			responseWithErrorMsg(w, http.StatusBadRequest, err)
			return
		}

		responseWithErrorMsg(w, http.StatusInternalServerError, err)
		return
	}

	responseWithJSON(w, http.StatusOK, u)
}
