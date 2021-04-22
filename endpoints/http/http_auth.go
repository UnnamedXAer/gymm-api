package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/unnamedxaer/gymm-api/usecases"
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
		setCookieJWTAuthToken(w, user.ID, app.jwtKey)
	}
	responseWithJSON(w, http.StatusOK, output)
}

func setCookieJWTAuthToken(w http.ResponseWriter, userID string, jwtKey []byte) {
	expirationTime := time.Now().Add(time.Second * 60)

	claims := Claims{
		ID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(jwtKey)
	if err != nil {
		responseWithInternalError(w)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cookieJwtTokenName,
		Value:    tokenStr,
		Expires:  expirationTime,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		Path:     "/",
	})
}

func deleteCookieJWTAuthToken(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   cookieJwtTokenName,
		MaxAge: -1, // remove cookie
	})
}
