package http

import (
	"context"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
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
			if err == http.ErrNoCookie {
				responseWithUnauthorized(w, "no token provided")
				return
			}
			responseWithError(w, http.StatusInternalServerError, err)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(cookie.Value, claims, func(t *jwt.Token) (interface{}, error) {
			return app.jwtKey, nil
		})
		if err != nil {
			responseWithUnauthorized(w, err)
			return
		}

		if claims.StandardClaims.ExpiresAt < time.Now().Add(30*time.Second).Unix() {
			device := r.UserAgent() // @todo: improve, mb send some info from client
			err := app.authUsecases.DeleteJWT(claims.UserID, device, cookie.Value)
			if err != nil {
				// we can ignore this error for because we are going to create new token anyway
				// if next calls fail we will return error to the client
				logDebugError(app.l, r, err)
			}

			rt, err := app.authUsecases.GetRefreshToken(claims.ID)
			if err != nil {
				logDebugError(app.l, r, err)
				responseWithInternalError(w)
				return
			}

			// if refresh token not exists, expired or is different then provided by client the user must login again

			err = validateRefreshToken(rt, claims.Token)
			if err != nil {
				app.authUsecases.DeleteRefreshToken(claims.ID)
				return
			}

			newToken, err := createJWTAuth(claims.ID, device, rt, app.jwtKey)
			if err != nil {
				logDebugError(app.l, r, err)
				responseWithInternalError(w)
				return
			}
			newToken, err = app.authUsecases.SaveJWT(
				newToken.UserID,
				newToken.Device,
				newToken.Token,
				newToken.ExpiresAt)
			if err != nil {
				logDebugError(app.l, r, err)
				responseWithInternalError(w)
				return
			}
			app.setCookieJWTAuthToken(w, newToken.Token, newToken.ExpiresAt)
		}

		if err = token.Claims.Valid(); err != nil {
			responseWithUnauthorized(w, err)
			return
		}
		ctx := context.WithValue(r.Context(), contextKeyUserID, claims.ID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}