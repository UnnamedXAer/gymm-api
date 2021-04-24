package http

import (
	"context"
	"net/http"

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
				responseWithUnauthorized(w, "no cookie")
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
		if err = token.Claims.Valid(); err != nil {
			responseWithUnauthorized(w, err)
			return
		}
		ctx := context.WithValue(r.Context(), contextKeyUserID, claims.ID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
