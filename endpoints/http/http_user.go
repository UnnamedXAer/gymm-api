package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/unnamedxaer/gymm-api/usecases"
	"github.com/unnamedxaer/gymm-api/validation"
)

func (app *App) CreateUser(w http.ResponseWriter, req *http.Request) {
	var u usecases.UserInput
	f, ff := validation.GetFieldJSONTag(&u, "Username")
	fmt.Println(f, ff)
	err := json.NewDecoder(req.Body).Decode(&u)
	log.Println("[POST / CreateUser] -> body: " + fmt.Sprintf("%v", u))
	if err != nil {
		resErrText := getErrOfMalformedInput(&u, []string{"ID", "CreatedAt"})
		responseWithErrorMsg(w, http.StatusUnprocessableEntity, errors.New(resErrText))
		return
	}
	defer req.Body.Close()

	err = validateUserInput(app.Validate, &u)
	if err != nil {
		if svErr, ok := err.(*validation.StructValidError); ok {
			responseWithErrorJSON(w, http.StatusNotAcceptable, svErr.ValidationErrors())
			return
		}
		responseWithErrorMsg(w, http.StatusInternalServerError, err)
		return
	}

	u.CreatedAt = time.Now()

	user, err := app.Usecases.CreateUser(&u)
	if err != nil {
		if strings.HasPrefix(err.Error(), "email address already in use") {
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
	log.Println("[GET / GetUserById] -> id: " + id)
	if id == "" {
		responseWithErrorMsg(w, http.StatusUnprocessableEntity, errors.New("Missign 'ID'"))
		return
	}

	u, err := app.Usecases.GetUserByID(id)
	if err != nil {
		responseWithErrorMsg(w, http.StatusInternalServerError, err)
		return
	}

	responseWithJSON(w, http.StatusOK, u)
}
