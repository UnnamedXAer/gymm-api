package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/unnamedxaer/gymm-api/entities"
	"github.com/unnamedxaer/gymm-api/repositories"
	"github.com/unnamedxaer/gymm-api/usecases"
	"github.com/unnamedxaer/gymm-api/validation"
)

type Claims struct {
	ID string `json:"id"`
	// IssuedAt is here to make the tokens unique by different payload
	IssuedAt               int64 `json:"createdAt"`
	*entities.RefreshToken `json:"r"`
	jwt.StandardClaims
}

const cookieJwtTokenName = "__HOST-token"

// Login validates input and login the user
func (app *App) Login(w http.ResponseWriter, req *http.Request) {
	var ui *usecases.UserInput
	err := json.NewDecoder(req.Body).Decode(&ui)
	if err != nil {
		logDebugError(app.l, req, err)

		resErrText := getErrOfMalformedInput(
			&ui, []string{"ID", "CreatedAt", "Username"})
		responseWithErrorTxt(w, http.StatusBadRequest, resErrText)
		return
	}

	ctx := req.Context()

	user, err := app.authUsecases.Login(ctx, ui)
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
	err = app.login(ctx, w, user.ID, device, req.Method, req.RequestURI)
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

	ctx := req.Context()

	user, err := app.userUsecases.CreateUser(ctx, &u)
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

	err = app.login(ctx, w, user.ID, req.UserAgent(), req.Method, req.RequestURI)
	if err != nil {
		logDebugError(app.l, req, err)
		clearCookieJWTAuthToken(w)
		responseWithInternalError(w)
		return
	}

	responseWithJSON(w, http.StatusCreated, output)
}

// Logout logouts the user, removes token from cookies & storage
func (app *App) Logout(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie(cookieJwtTokenName)
	var n int64
	if cookie != nil {
		ctx := req.Context()
		n, err = app.authUsecases.DeleteJWT(
			ctx,
			&entities.UserToken{
				Token: cookie.Value,
			})

	}
	if err != nil {
		logDebugError(app.l, req, err)
	} else if n == 0 {
		logDebugError(app.l, req, fmt.Errorf(
			"JWT was not deleted for cookie.value: %s", cookie.Value))
	}
	clearCookieJWTAuthToken(w)
	w.WriteHeader(http.StatusOK)
}

// GetSessions returns all user sessions
func (app *App) GetSessions(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userID, ok := ctx.Value(contextKeyUserID).(string)
	if !ok {
		responseWithUnauthorized(w)
		return
	}

	jwts, err := app.authUsecases.GetUserJWTs(ctx, userID, entities.All)
	if err != nil {
		logDebugError(app.l, req, err)
		clearCookieJWTAuthToken(w)
		responseWithInternalError(w)
		return
	}

	responseWithJSON(w, http.StatusOK, &jwts)
}

// LogoutSession logouts the user, from the device/browser
func (app *App) LogoutSession(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userID, ok := ctx.Value(contextKeyUserID).(string)
	if !ok {
		responseWithUnauthorized(w)
		return
	}

	ut := entities.UserToken{}
	err := json.NewDecoder(req.Body).Decode(&ut)
	if err != nil {
		logDebugError(app.l, req, err)
		responseWithError(w, http.StatusBadRequest, err)
		return
	}
	if ut.UserID != "" && ut.UserID != userID {
		responseWithUnauthorized(w)
		return
	}

	ut.UserID = userID

	cookie, _ := req.Cookie(cookieJwtTokenName)
	if cookie != nil {
		// determine if user want to delete current token
		// and logout his in that case
		if ut.Token != "" {
			if ut.Token == cookie.Value {
				clearCookieJWTAuthToken(w)
			}
		} else {
			storedTokens, err := app.authUsecases.GetUserJWTs(
				ctx, userID, entities.NotExpired)
			if err != nil {
				logDebugError(app.l, req, err)
			}

			for _, t := range storedTokens {
				if t.Token == cookie.Value {
					clearCookieJWTAuthToken(w)
					break
				}
			}
		}
	}

	n, err := app.authUsecases.DeleteJWT(ctx, &ut)
	if err != nil {
		logDebugError(app.l, req, err)
		var e *repositories.InvalidIDError
		if errors.As(err, &e) {
			responseWithError(w, http.StatusBadRequest, e)
			return
		}
		responseWithInternalError(w)
		return
	}

	if n == 0 {
		responseWithJSON(
			w,
			http.StatusOK,
			map[string]string{"warning": "no records were deleted"})
		return
	}

	w.WriteHeader(http.StatusOK)
}

// LogoutSession logouts the user, from the device/browser
func (app *App) LogoutAllSessions(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userID, ok := ctx.Value(contextKeyUserID).(string)
	if !ok {
		responseWithUnauthorized(w)
		return
	}

	n, err := app.authUsecases.DeleteRefreshTokenAndAllTokens(ctx, userID)
	if err != nil {
		logDebugError(app.l, req, err)
		var e *repositories.InvalidIDError
		if errors.As(err, &e) {
			responseWithError(w, http.StatusBadRequest, e)
			return
		}
		responseWithInternalError(w)
		return
	}
	clearCookieJWTAuthToken(w)
	if n == 0 {
		responseWithJSON(
			w, http.StatusOK, map[string]string{"warning": "no records were deleted"})
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (app *App) ChangePassword(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userID, ok := ctx.Value(contextKeyUserID).(string)
	if !ok {
		responseWithUnauthorized(w)
		return
	}

	body := make(map[string]string, 2)
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		logDebugError(app.l, req, err)
		responseWithError(w, http.StatusBadRequest, err)
		return
	}

	rv := reflect.TypeOf(usecases.UserInput{})
	srf, ok := rv.FieldByName("Password")
	if !ok {
		responseWithInternalError(w)
		return
	}

	validationRules := srf.Tag.Get("validate")
	fieldName := "password"
	err = app.Validate.Var(body[fieldName], validationRules)
	if err != nil {
		logDebugError(app.l, req, err)
		// @refactor: see exercise validation
		formattedErrors := make(map[string]string, 1)
		validateErrs, ok := err.(validator.ValidationErrors)
		if !ok {
			formattedErrors[fieldName] = err.Error()
		}

		for _, err := range validateErrs {
			formattedErrors[fieldName] += getErrorTranslation4User(&err, fieldName)
		}

		vErrs := validation.NewStructValidError(formattedErrors)
		responseWithJSON(w, http.StatusBadRequest, vErrs.Format())
		return
	}

	err = app.authUsecases.ChangePassword(ctx, userID, body["oldPassword"], body["password"])
	if err != nil {
		logDebugError(app.l, req, err)
		if errors.Is(err, usecases.IncorrectCredentialsError{}) {
			responseWithError(w, http.StatusUnauthorized, err)
			return
		}
		// @todo: handle other error types
		responseWithError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *App) Refresh(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// login make user authenticated by
// creating tokens and setting them in user response as a cookie
// tokens are also stored in db
func (app *App) login(
	ctx context.Context,
	w http.ResponseWriter,
	userID, device, method, reqURI string) error {
	rt, err := app.authUsecases.GetRefreshToken(ctx, userID)
	if err != nil {
		app.l.Debug().Msgf("[%s %s]: error: %v", method, reqURI, err)
		// the retrive of the refresh token has failed for some reason,
		// but we can try to create new refresh token.
	}

	if rt == nil || time.Until(rt.ExpiresAt) <= 0 {
		rt, err = createRefreshToken(
			ctx, userID, app.authUsecases.SaveRefreshToken)
		if err != nil {
			return err
		}
	}

	token, err := createJWTAuth(
		ctx, userID, device, rt, app.jwtKey, app.authUsecases.SaveJWT)
	if err != nil {
		return err
	}

	setCookieJWTAuthToken(w, token.Token, token.ExpiresAt)
	return nil
}

// sets a cookie with jwt
func setCookieJWTAuthToken(w http.ResponseWriter, token string, expTime time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:  cookieJwtTokenName,
		Value: token,
		// Expires:  expTime,
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
	ctx context.Context,
	userID string,
	device string,
	rt *entities.RefreshToken,
	jwtKey []byte,
	saveFunc func(
		ctx context.Context,
		userID string,
		device string,
		token string,
		expiresAt time.Time) (*entities.UserToken, error)) (*entities.UserToken, error) {

	now := time.Now().UTC()
	expirationTime := now.Add(time.Second * 60 * 5) // @todo: configurable time

	claims := Claims{
		ID:           userID,
		RefreshToken: rt,
		IssuedAt:     now.UnixNano(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  now.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	tokenStr, err := token.SignedString(jwtKey)
	if err != nil {
		return nil, err
	}

	return saveFunc(ctx, userID, device, tokenStr, expirationTime)
}

// generate new refresh token for given user and saves it in storage
func createRefreshToken(
	ctx context.Context,
	userID string,
	saveFunc func(
		ctx context.Context,
		userID string,
		token string,
		expiresAt time.Time) (*entities.RefreshToken, error)) (*entities.RefreshToken, error) {

	refreshVal, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return saveFunc(ctx, userID, refreshVal.String(), time.Now().AddDate(1, 0, 0))
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
