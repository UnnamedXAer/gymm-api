package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/unnamedxaer/gymm-api/usecases"
)

func CreateUser(w http.ResponseWriter, req *http.Request) {
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
	// logger := zerolog.Logger{}
	// db, err := repositories.GetDatabase(&logger, os.Getenv("MONGO_URI"), os.Getenv("DB_NAME"))
	// createUser := usecases.CreateUserUseCase(usersRepo.NewRepository(
	// 	&logger,
	// 	db.Collection("users"),
	// ))

	user, err := createUser(&u)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, err)
		return
	}
	// u.Password = ""
	responseWithJSON(w, http.StatusCreated, user)
}

// func GetUserById(w http.ResponseWriter, req *http.Request) {

// 	vars := mux.Vars(req)
// 	id := vars["id"]
// 	log.Println("[GET / GetUserById] -> id: " + id)
// 	var u models.User
// 	if id == "" {
// 		responseWithError(w, http.StatusUnprocessableEntity, errors.New("Missign 'ID'"))
// 		return
// 	}
// 	u.ID = id

// 	err := services.UService.GetUserById(&u)
// 	if err != nil {
// 		responseWithError(w, http.StatusInternalServerError, err)
// 		return
// 	}

// 	responseWithJSON(w, http.StatusOK, u)
// }
