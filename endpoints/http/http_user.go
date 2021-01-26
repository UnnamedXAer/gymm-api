package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/unnamedxaer/gymm-api/junk/server/models"
	"github.com/unnamedxaer/gymm-api/junk/services"
)

func CreateUser(w http.ResponseWriter, req *http.Request) {
	var u models.User
	err := json.NewDecoder(req.Body).Decode(&u)
	if err != nil {
		responseWithError(w, http.StatusUnprocessableEntity, err)
		return
	}
	defer req.Body.Close()
	u.CreatedAt = time.Now()
	log.Println("[CreateUser] -> body: " + fmt.Sprintf("%v", u))
	err = services.UService.CreateUser(&u)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, err)
		return
	}
	// u.Password = ""
	responseWithJSON(w, http.StatusCreated, u)
}

func GetUserById(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	id := vars["id"]
	log.Println("[GetUserById] -> id: " + id)
	var u models.User
	if id == "" {
		responseWithError(w, http.StatusUnprocessableEntity, errors.New("Missign 'ID'"))
		return
	}
	u.ID = id

	err := services.UService.GetUserById(&u)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, err)
		return
	}

	responseWithJSON(w, http.StatusOK, u)
}
