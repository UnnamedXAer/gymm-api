package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/unnamedxaer/gymm-api/usecases"
)

func (app *App) CreateUser(w http.ResponseWriter, req *http.Request) {
	var u usecases.UserData
	err := json.NewDecoder(req.Body).Decode(&u)
	if err != nil {
		responseWithError(w, http.StatusUnprocessableEntity, err)
		return
	}
	defer req.Body.Close()
	u.CreatedAt = time.Now()
	log.Println("[POST / CreateUser] -> body: " + fmt.Sprintf("%v", u))

	// @todo: refactor
	panic("not implemented yet")

	user, err := app.Usecases.CreateUser(&u)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, err)
		return
	}
	// u.Password = ""
	responseWithJSON(w, http.StatusCreated, user)
}

func (app *App) GetUserById(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	id := vars["id"]
	log.Println("[GET / GetUserById] -> id: " + id)
	if id == "" {
		responseWithError(w, http.StatusUnprocessableEntity, errors.New("Missign 'ID'"))
		return
	}

	u, err := app.Usecases.GetUserByID(id)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, err)
		return
	}

	responseWithJSON(w, http.StatusOK, u)
}
