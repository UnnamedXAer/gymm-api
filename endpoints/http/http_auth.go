package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/unnamedxaer/gymm-api/repositories"
	"github.com/unnamedxaer/gymm-api/usecases"
	"github.com/unnamedxaer/gymm-api/validation"
)

type Claims struct {
	ID string `json:"id"`
	jwt.StandardClaims
}

const cookieJwtTokenName = "__Host-token"

func (app *App) Login(w http.ResponseWriter, req *http.Request) {
	var ui *usecases.UserInput
	err := json.NewDecoder(req.Body).Decode(&ui)
	if err != nil {
		logDebugError(app.l, req, err)

		resErrText := getErrOfMalformedInput(&ui, []string{"ID", "CreatedAt", "Username"})
		responseWithErrorTxt(w, http.StatusBadRequest, resErrText)
		return
	}

	user, err := app.authUsecases.Login(ui)
	if err != nil && errors.Is(err, &usecases.IncorrectCredentialsError{}) {
		logDebugError(app.l, req, err)
		responseWithInternalError(w)
		return
	}

	output := map[string]interface{}{
		"user": user,
	}

	if user == nil {
		output["error"] = "incorrect credentials"
		deleteCookieJWTAuthToken(w)
	} else {
		err := setCookieJWTAuthToken(w, user.ID, app.jwtKey)
		if err != nil {
			logDebugError(app.l, req, err)
			responseWithInternalError(w)
		}
	}
	responseWithJSON(w, http.StatusOK, output)
}

func (app *App) Register(w http.ResponseWriter, req *http.Request) {
	var u usecases.UserInput

	err := json.NewDecoder(req.Body).Decode(&u)
	app.l.Debug().Msg("[POST / CreateUser] -> body: " + fmt.Sprintf("%v", u))
	if err != nil {
		resErrText := getErrOfMalformedInput(&u, []string{"ID", "CreatedAt"})
		responseWithError(w, http.StatusUnprocessableEntity, errors.New(resErrText))
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
		responseWithError(w, http.StatusInternalServerError, err)
		return
	}

	user, err := app.userUsecases.CreateUser(&u)
	if err != nil {
		if errors.Is(err, repositories.NewErrorEmailAddressInUse()) {
			responseWithError(w, http.StatusConflict, err)
			return
		}

		responseWithError(w, http.StatusInternalServerError, err)
		return
	}

	responseWithJSON(w, http.StatusCreated, user)
}

func (app *App) Logout(w http.ResponseWriter, req *http.Request) {
	deleteCookieJWTAuthToken(w)
	w.WriteHeader(http.StatusOK)
}

func setCookieJWTAuthToken(w http.ResponseWriter, userID string, jwtKey []byte) error {
	expirationTime := time.Now().Add(time.Second * 60)

	claims := Claims{
		ID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	tokenStr, err := token.SignedString(jwtKey)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cookieJwtTokenName,
		Value:    tokenStr,
		Expires:  expirationTime,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		// Secure:   true,
		Path: "/",
	})

	return nil
}

func deleteCookieJWTAuthToken(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   cookieJwtTokenName,
		MaxAge: -1, // remove cookie
	})
}
