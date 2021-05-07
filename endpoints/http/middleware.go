package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/unnamedxaer/gymm-api/entities"
)

type middleware func(http.HandlerFunc) http.HandlerFunc

func chainMiddlewares(h http.HandlerFunc, middlewares ...middleware) http.HandlerFunc {
	if len(middlewares) == 0 {
		return h
	}

	wrapped := h

	for i := len(middlewares) - 1; i >= 0; i-- {
		wrapped = middlewares[i](wrapped)
	}

	return wrapped
}

func (app *App) checkAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(cookieJwtTokenName)
		if err != nil {
			logDebugError(app.l, r, err)
			if err == http.ErrNoCookie {
				responseWithUnauthorized(w, "no token provided")
				return
			}
			responseWithError(w, http.StatusInternalServerError, err)
			return
		}

		cookieVal := cookie.Value
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(
			cookieVal, claims,
			func(t *jwt.Token) (interface{}, error) {
				return app.jwtKey, nil
			})
		if err != nil {
			var vErr *jwt.ValidationError
			if !(errors.As(err, &vErr) &&
				// verify that only expiration time is not valid
				// @thought: maybe validate time manually eg. by using IssuedAt + "max age"
				// instead of expiration time, this way there will be no need of asserting
				// the error (if exptime will be zero value)
				(vErr.Errors == jwt.ValidationErrorExpired)) ||
				claims.ID == "" ||
				claims.StandardClaims.ExpiresAt == 0 {
				logDebugError(app.l, r, err)
				clearCookieJWTAuthToken(w)
				responseWithUnauthorized(w, err)
				return
			}
		}

		ctx := r.Context()

		if claims.StandardClaims.ExpiresAt < time.Now().Add(30*time.Second).Unix() {
			device := r.UserAgent() // @todo: improve, mb send some info from client
			ut := entities.UserToken{
				UserID: claims.UserID,
				Token:  cookie.Value,
				Device: device,
			}
			n, err := app.authUsecases.DeleteJWT(ctx, &ut)
			if err != nil {
				// we can ignore this error for because we are going to create new token anyway
				// if next calls fail we will return error to the client
				logDebugError(app.l, r, err)
			} else if n == 0 {
				logDebugError(app.l, r, fmt.Errorf("JWT was not deleted for: %v", ut))
			}

			rt, err := app.authUsecases.GetRefreshToken(ctx, claims.ID)
			if err != nil {
				logDebugError(app.l, r, err)
				clearCookieJWTAuthToken(w)
				responseWithInternalError(w)
				return
			}

			// if refresh token not exists, expired or is different than
			// provided by client the user must login again

			err = validateRefreshToken(rt, claims.Token)
			if err != nil {
				logDebugError(app.l, r, err)
				// @todo: handle return
				_, err = app.authUsecases.DeleteRefreshToken(ctx, claims.ID)
				if err != nil {
					logDebugError(app.l, r, err)
				}
				clearCookieJWTAuthToken(w)
				responseWithUnauthorized(w, err)
				return
			}

			newToken, err := createJWTAuth(ctx, claims.ID, device, rt, app.jwtKey, app.authUsecases.SaveJWT)
			if err != nil {
				logDebugError(app.l, r, err)
				responseWithInternalError(w)
				return
			}
			setCookieJWTAuthToken(w, newToken.Token, newToken.ExpiresAt)
		} else {
			if err = token.Claims.Valid(); err != nil {
				clearCookieJWTAuthToken(w)
				responseWithUnauthorized(w, err)
				return
			}
		}

		ctx = context.WithValue(ctx, contextKeyUserID, claims.ID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func suffixMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}
		next.ServeHTTP(w, r)
	})
}
