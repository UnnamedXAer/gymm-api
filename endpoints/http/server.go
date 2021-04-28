package http

import (
	"net/http"

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
	router.StrictSlash(true)

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
	app.Router.HandleFunc("/health", chainMiddlewares(app.Health, app.checkAuthenticated)).Methods(http.MethodGet)

	// app.Router.HandleFunc(
	// 	"/users/{id:[0-9a-zA-Z]+}",
	// 	chainMiddlewares(app.GetUserById, app.checkAuthenticated)).Methods(http.MethodGet)

	exercisesRouter := app.Router.PathPrefix("/exercises").Subrouter()
	exercisesRouter.HandleFunc(
		"/{exerciseID:[0-9a-zA-Z]+}",
		chainMiddlewares(app.GetExerciseByID, app.checkAuthenticated)).Methods(http.MethodGet)
	exercisesRouter.HandleFunc(
		"",
		chainMiddlewares(app.GetExercisesByName, app.checkAuthenticated)).Methods(http.MethodGet)
	exercisesRouter.HandleFunc(
		"/{exerciseID:[0-9a-zA-Z]+}",
		chainMiddlewares(app.UpdateExercise, app.checkAuthenticated)).Methods(http.MethodPatch)
	exercisesRouter.HandleFunc(
		"",
		chainMiddlewares(app.CreateExercise, app.checkAuthenticated)).Methods(http.MethodPost)

	// training
	trainingRouter := app.Router.PathPrefix("/trainings").Subrouter()
	trainingRouter.HandleFunc(
		"",
		chainMiddlewares(app.GetUserTrainings, app.checkAuthenticated)).Methods(http.MethodGet)
	trainingRouter.HandleFunc(
		"",
		chainMiddlewares(app.StartTraining, app.checkAuthenticated)).Methods(http.MethodPost)
	trainingRouter.HandleFunc(
		"/{trainingID:[0-9a-zA-Z]+}",
		chainMiddlewares(app.GetTrainingByID, app.checkAuthenticated)).Methods(http.MethodGet)
	trainingRouter.HandleFunc(
		"/{trainingID:[0-9a-zA-Z]+}/end",
		chainMiddlewares(app.EndTraining, app.checkAuthenticated)).Methods(http.MethodPatch)

	// training exercise
	trainingExerciseRouter := trainingRouter.PathPrefix("/{trainingID:[0-9a-zA-Z]+}/exercises").Subrouter()
	trainingExerciseRouter.HandleFunc(
		"",
		chainMiddlewares(app.StartTrainingExercise, app.checkAuthenticated)).Methods(http.MethodPost)
	trainingExerciseRouter.HandleFunc(
		"/{exerciseID:[0-9a-zA-Z]+}/end",
		chainMiddlewares(app.EndTrainingExercise, app.checkAuthenticated)).Methods(http.MethodPatch)

	// training set
	trainingSetRouter := trainingExerciseRouter.PathPrefix("/{exerciseID:[0-9a-zA-Z]+}/sets").Subrouter()
	trainingSetRouter.HandleFunc(
		"",
		chainMiddlewares(app.AddTrainingSetExercise, app.checkAuthenticated)).Methods(http.MethodPost)

	app.Router.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		logDebug(app.l, r, nil)
		rw.WriteHeader(http.StatusMethodNotAllowed)
	})
}

func (app *App) Run(addr string) {
	srv := &http.Server{
		Addr:    addr,
		Handler: suffixMiddleware(app.Router),
	}
	app.l.Info().Msg("server is up and running at " + addr)
	app.l.Error().Stack().Err(srv.ListenAndServe()).Msg("")
}
