package http

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/unnamedxaer/gymm-api/usecases"
)

type App struct {
	l        *zerolog.Logger
	Usecases usecases.IUserUseCases
	Router   *mux.Router
	Validate *validator.Validate
}

func NewServer(
	logger *zerolog.Logger,
	userRepo usecases.UserRepo,
	validate *validator.Validate) *App {
	userUsecases := usecases.NewUserUseCases(userRepo)
	router := mux.NewRouter()
	app := App{
		l:        logger,
		Usecases: userUsecases,
		Router:   router,
		Validate: validate,
	}
	return &app
}

func (app *App) AddHandlers() {
	app.Router.HandleFunc("/users/{id:[0-9a-zA-Z]+}", app.GetUserById).Methods("GET")
	app.Router.HandleFunc("/users", app.CreateUser).Methods("POST")

	app.Router.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		app.l.Info().Msg("[" + r.Method + "/] -> URL: " + r.RequestURI)
		rw.WriteHeader(http.StatusMethodNotAllowed)
	})
}

func (app *App) Run(addr string) {
	app.l.Info().Msg("server is up and running at " + addr)
	app.l.Error().Stack().Err(http.ListenAndServe(addr, app.Router)).Msg("")
}
