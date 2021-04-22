package http

import (
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/unnamedxaer/gymm-api/usecases"
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
	validate *validator.Validate) *App {

	jwtKey := []byte(os.Getenv("JWT_KEY"))
	if len(jwtKey) < 10 {
		logger.Panic().Msg("missing or too short jwt key")
	}

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

	app.Router.HandleFunc("/login", app.Login).Methods("POST")

	app.Router.HandleFunc("/users/{id:[0-9a-zA-Z]+}", app.GetUserById).Methods("GET")
	app.Router.HandleFunc("/users", app.CreateUser).Methods("POST")

	app.Router.HandleFunc("/exercises/{id:[0-9a-zA-Z]+}", app.GetExeriseByID).Methods(http.MethodGet)
	app.Router.HandleFunc("/exercises/{id:[0-9a-zA-Z]+}", app.UpdateExercise).Methods(http.MethodPatch)
	app.Router.HandleFunc("/exercises", app.CreateExercise).Methods(http.MethodPost)

	app.Router.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		app.l.Debug().Msg("[" + r.Method + "/] -> URL: " + r.RequestURI)
		rw.WriteHeader(http.StatusMethodNotAllowed)
	})
}

func (app *App) Run(addr string) {
	app.l.Info().Msg("server is up and running at " + addr)
	app.l.Error().Stack().Err(http.ListenAndServe(addr, app.Router)).Msg("")
}
