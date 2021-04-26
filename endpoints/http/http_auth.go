package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/unnamedxaer/gymm-api/entities"
	"github.com/unnamedxaer/gymm-api/repositories"
	"github.com/unnamedxaer/gymm-api/usecases"
	"github.com/unnamedxaer/gymm-api/validation"
)

type Claims struct {
	ID                     string `json:"id"`
	*entities.RefreshToken `json:"r"`
	jwt.StandardClaims
}

const cookieJwtTokenName = "__Host-token"

// Login validates input and login the user
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
		clearCookieJWTAuthToken(w)
		responseWithJSON(w, http.StatusOK, output)
		return
	}

	device := req.UserAgent() // @todo: improve, mb send some info from client

	err = app.login(w, user.ID, device, req.Method, req.RequestURI)
	if err != nil {
		logDebugError(app.l, req, err)
		clearCookieJWTAuthToken(w)
		responseWithInternalError(w)
		return
	}

	responseWithJSON(w, http.StatusOK, output)
}

// Register validates input, creates new user and login the user
func (app *App) Register(w http.ResponseWriter, req *http.Request) {
	var u usecases.UserInput

	err := json.NewDecoder(req.Body).Decode(&u)
	if err != nil {
		logDebugError(app.l, req, err)
		resErrText := getErrOfMalformedInput(&u, []string{"ID", "CreatedAt"})
		responseWithError(w, http.StatusUnprocessableEntity, errors.New(resErrText))
		return
	}
	defer req.Body.Close()

	trimWhitespacesOnUserInput(&u)

	err = validateUserInput(app.Validate, &u)
	if err != nil {
		logDebugError(app.l, req, err)
		if svErr, ok := err.(*validation.StructValidError); ok {
			responseWithJSON(w, http.StatusNotAcceptable, svErr.ValidationErrors())
			return
		}
		responseWithError(w, http.StatusInternalServerError, err)
		return
	}

	user, err := app.userUsecases.CreateUser(&u)
	if err != nil {
		logDebugError(app.l, req, err)
		if errors.Is(err, repositories.NewErrorEmailAddressInUse()) {
			responseWithErrorTxt(w, http.StatusConflict, "email address already in use")
			return
		}

		responseWithError(w, http.StatusInternalServerError, err)
		return
	}

	// login user

	output := map[string]interface{}{
		"user": user,
	}

	err = app.login(w, user.ID, req.UserAgent(), req.Method, req.RequestURI)
	if err != nil {
		logDebugError(app.l, req, err)
		clearCookieJWTAuthToken(w)
		responseWithInternalError(w)
		return
	}

	responseWithJSON(w, http.StatusOK, output)
}

// Logout logouts the user, removes token from cookies & storage
func (app *App) Logout(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie(cookieJwtTokenName)
	if cookie != nil {
		err = app.authUsecases.DeleteJWT("", "", cookie.Value)
	}
	if err != nil {
		logDebugError(app.l, req, err)
	}
	clearCookieJWTAuthToken(w)
	w.WriteHeader(http.StatusOK)
}

func (app *App) Refresh(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// login make user authenticated by
// creating tokens and setting them in user response as a cookie
// tokens are also stored in db
func (app *App) login(w http.ResponseWriter, userID, device, method, reqURI string) error {
	rt, err := app.authUsecases.GetRefreshToken(userID)
	if err != nil {
		app.l.Debug().Msgf("[%s %s]: error: %v", method, reqURI, err)
		// the retrive of the refresh token has failed for some reason,
		// but we can try to create new refresh token.
	}

	if rt == nil || time.Until(rt.ExpiresAt) <= 0 {
		rt, err = createRefreshToken(userID, app.authUsecases.SaveRefreshToken)
		if err != nil {
			return err
		}
	}

	token, err := createJWTAuth(userID, device, rt, app.jwtKey, app.authUsecases.SaveJWT)
	if err != nil {
		return err
	}

	setCookieJWTAuthToken(w, token.Token, token.ExpiresAt)
	return nil
}

// sets a cookie with jwt
func setCookieJWTAuthToken(w http.ResponseWriter, token string, expTime time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieJwtTokenName,
		Value:    token,
		Expires:  expTime,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		// Secure:   true,
		Path: "/",
	})
}

// sets -1 to the cookie max age to force its removal
func clearCookieJWTAuthToken(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   cookieJwtTokenName,
		MaxAge: -1, // remove cookie
	})
}

// generate new jwt for given user and saves it in storage
func createJWTAuth(
	userID string,
	device string,
	rt *entities.RefreshToken,
	jwtKey []byte,
	saveFunc func(
		userID string,
		device string,
		token string,
		expiresAt time.Time) (*entities.UserToken, error)) (*entities.UserToken, error) {

	expirationTime := time.Now().Add(time.Second * 60 * 5) // @todo: configurable time

	claims := Claims{
		ID:           userID,
		RefreshToken: rt,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	tokenStr, err := token.SignedString(jwtKey)
	if err != nil {
		return nil, err
	}

	return saveFunc(userID, device, tokenStr, expirationTime)
}

// generate new refresh token for given user and saves it in storage
func createRefreshToken(
	userID string,
	saveFunc func(
		userID string,
		token string,
		expiresAt time.Time) (*entities.RefreshToken, error)) (*entities.RefreshToken, error) {

	refreshVal, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return saveFunc(userID, refreshVal.String(), time.Now().AddDate(1, 0, 0))
}

// it returns nil error if token is valid
func validateRefreshToken(rt *entities.RefreshToken, clientRefToken string) error {
	if rt == nil || time.Until(rt.ExpiresAt) <= 0 {
		return errors.New("token & refresh token expired")
	}
	if rt.Token != clientRefToken {
		return errors.New("invalid refresh token")
	}
	return nil
}
