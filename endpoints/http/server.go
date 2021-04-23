package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/unnamedxaer/gymm-api/usecases"
)

type contextKey int

const (
	contextKeyUserID contextKey = iota
)

type App struct {
	l                *zerolog.Logger
	authUsecases     usecases.IAuthUsecases
	userUsecases     usecases.IUserUseCases
	exerciseUsecases usecases.IExerciseUseCases
	trainingUsecases usecases.ITrainingUsecases
	Router           *mux.Router
	Validate         *validator.Validate
	jwtKey           []byte
}

func NewServer(
	logger *zerolog.Logger,
	authRepo usecases.AuthRepo,
	userRepo usecases.UserRepo,
	exerciseRepo usecases.ExerciseRepo,
	trainingRepo usecases.TrainingRepo,
	validate *validator.Validate,
	jwtKey []byte) *App {

	var authUsecases usecases.IAuthUsecases = usecases.NewAuthUsecases(authRepo)
	var userUsecases usecases.IUserUseCases = usecases.NewUserUseCases(userRepo)
	var exerciseUsecases usecases.IExerciseUseCases = usecases.NewExerciseUseCases(exerciseRepo)
	var trainingUsecases usecases.ITrainingUsecases = usecases.NewTrainingUseCases(trainingRepo)

	router := mux.NewRouter()
	app := App{
		l:                logger,
		authUsecases:     authUsecases,
		userUsecases:     userUsecases,
		exerciseUsecases: exerciseUsecases,
		trainingUsecases: trainingUsecases,
		Router:           router,
		Validate:         validate,
		jwtKey:           jwtKey,
	}
	return &app
}

func (app *App) AddHandlers() {

	app.Router.HandleFunc("/login", app.Login).Methods(http.MethodPost)
	app.Router.HandleFunc("/logout", app.Logout).Methods(http.MethodGet)
	app.Router.HandleFunc("/register", app.Register).Methods(http.MethodPost)

	app.Router.HandleFunc("/users/{id:[0-9a-zA-Z]+}", app.GetUserById).Methods("GET")

	app.Router.HandleFunc("/exercises/{id:[0-9a-zA-Z]+}", app.GetExeriseByID).Methods(http.MethodGet)
	app.Router.HandleFunc("/exercises/{id:[0-9a-zA-Z]+}", app.UpdateExercise).Methods(http.MethodPatch)
	// app.Router.HandleFunc("/exercises", app.CreateExercise).Methods(http.MethodPost)
	app.Router.Handle("/exercises",
		app.loggingMdw(app.checkAuthenticated(
			http.HandlerFunc(app.CreateExercise)))).Methods(http.MethodPost)

	// app.Router.Use()
	app.Router.HandleFunc("/heath", app.Health).Methods(http.MethodGet)

	app.Router.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		logDebug(app.l, r, nil)
		rw.WriteHeader(http.StatusMethodNotAllowed)
	})
}

func (app *App) Run(addr string) {
	app.l.Info().Msg("server is up and running at " + addr)
	app.l.Error().Stack().Err(http.ListenAndServe(addr, app.Router)).Msg("")
}

func (app *App) checkAuthenticated(next http.Handler) http.Handler {
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

func (app *App) loggingMdw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println()
		fmt.Println("loggingMdw -> " + r.RequestURI)
		fmt.Println()
		next.ServeHTTP(w, r)
	})
}
