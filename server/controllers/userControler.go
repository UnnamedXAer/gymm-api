package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/unnamedxaer/gymm-api/server/models"
	"github.com/unnamedxaer/gymm-api/services"
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
	id1 := req.FormValue("id")
	log.Printf("vars: %v", vars)
	id := vars["id"]
	if id == "" {
		id = id1
	}
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
