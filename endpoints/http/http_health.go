package http

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func (app *App) Health(w http.ResponseWriter, req *http.Request) {
	output := make(map[string]interface{})

	defer logDebug(app.l, req, output)

	cookie, err := req.Cookie(cookieJwtTokenName)
	if err != nil {
		if err == http.ErrNoCookie {
			output["token"] = "no cookie"
		} else {
			output["token"] = "could not retrieve"
		}
		output["error"] = err.Error()
		responseWithJSON(w, http.StatusUnauthorized, &output)
		return
	}

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(t *jwt.Token) (interface{}, error) {
		return app.jwtKey, nil
	})
	if err != nil {
		output["token"] = "corrupted"
		output["error"] = err.Error()
		responseWithJSON(w, http.StatusUnauthorized, &output)
		return
	}
	if err = token.Claims.Valid(); err != nil {
		output["token"] = "corrupted claims"
		output["error"] = err.Error()
		responseWithJSON(w, http.StatusUnauthorized, &output)
		return
	}

	if claims.StandardClaims.ExpiresAt <= time.Now().Unix() {
		output["token"] = "expired"
		output["error"] = "token expired"
		responseWithJSON(w, http.StatusUnauthorized, &output)
		return
	}

	output["token"] = "OK"

	loggedUserID := claims.ID

	ctx := req.Context()

	user, err := app.userUsecases.GetUserByID(ctx, loggedUserID)
	if err != nil {
		output["DB"] = "error"
		output["user"] = "could not retrieve"
		output["error"] = err.Error()
		responseWithJSON(w, http.StatusUnauthorized, &output)
		return
	}
	output["DB"] = "OK"

	if user == nil {
		output["user"] = "user not found"
		output["error"] = "user not found"
		responseWithJSON(w, http.StatusUnauthorized, &output)
		return
	}

	output["user"] = user
	output["error"] = nil
	responseWithJSON(w, http.StatusOK, &output)

}
