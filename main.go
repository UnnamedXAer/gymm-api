package main

import (
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	"github.com/unnamedxaer/gymm-api/endpoints/http"
	"github.com/unnamedxaer/gymm-api/mailer"
	"github.com/unnamedxaer/gymm-api/repositories"
	"github.com/unnamedxaer/gymm-api/repositories/auth"
	"github.com/unnamedxaer/gymm-api/repositories/exercises"
	"github.com/unnamedxaer/gymm-api/repositories/trainings"
	"github.com/unnamedxaer/gymm-api/repositories/users"
	"github.com/unnamedxaer/gymm-api/validation"
)

func main() {
	logger := zerolog.New(os.Stdout)
	logger.Info().Msg(time.Now().Local().String() + "-> App starts, env = " + os.Getenv("ENV"))

	// @refactor: make it config file
	port := os.Getenv("PORT")
	if port == "" {
		panic("environment variable 'PORT' is not set")
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		panic("environment variable 'DB_NAME' is not set")
	}
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		panic("environment variable 'MONGO_URI' is not set")
	}
	jwtKey := []byte(os.Getenv("JWT_KEY"))
	if len(jwtKey) < 10 {
		panic("environment variable 'JWT_KEY' is not set or is too short")
	}

	db, err := repositories.GetDatabase(&logger, mongoURI, dbName)
	if err != nil {
		panic(err)
	}
	err = repositories.CreateCollections(&logger, db)
	if err != nil {
		logger.Panic().Msg(err.Error())
	}

	usersCol := repositories.GetCollection(&logger, db, repositories.UsersCollectionName)
	tokensCol := repositories.GetCollection(&logger, db, repositories.TokensCollectionName)
	refTokensCol := repositories.GetCollection(&logger, db, repositories.RefreshTokensCollectionName)
	resPwdReqsCol := repositories.GetCollection(&logger, db, repositories.ResPwdReqCollectionName)
	usersRepo := users.NewRepository(&logger, usersCol)

	authRepo := auth.NewRepository(&logger, usersCol, tokensCol, refTokensCol, resPwdReqsCol)

	exercisesCol := repositories.GetCollection(&logger, db, repositories.ExercisesCollectionName)
	exercisesRepo := exercises.NewRepository(&logger, exercisesCol)

	trainingsCol := repositories.GetCollection(&logger, db, repositories.TrainingsCollectionName)
	trainingsRepo := trainings.NewRepository(&logger, trainingsCol)

	validate := validation.New()

	mailer := mailer.NewMailer(&logger, func(err error) {
		logger.Err(err).Send()
	})

	app := http.NewServer(
		&logger,
		authRepo,
		usersRepo,
		exercisesRepo,
		trainingsRepo,
		validate,
		jwtKey,
		mailer,
	)

	app.AddHandlers()

	app.Run("localhost:" + os.Getenv("PORT"))
}
