package main

import (
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	"github.com/unnamedxaer/gymm-api/endpoints/http"
	"github.com/unnamedxaer/gymm-api/repositories"
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

	db, err := repositories.GetDatabase(&logger, mongoURI, dbName)
	if err != nil {
		panic(err)
	}
	repositories.CreateCollections(&logger, db)

	usersCol := repositories.GetCollection(&logger, db, repositories.UsersCollectionName)
	usersRepo := users.NewRepository(&logger, usersCol)

	exercisesCol := repositories.GetCollection(&logger, db, repositories.ExercisesCollectionName)
	exercisesRepo := exercises.NewRepository(&logger, exercisesCol)

	trainingsCol := repositories.GetCollection(&logger, db, repositories.TrainingsCollectionName)
	trainingsRepo := trainings.NewRepository(&logger, trainingsCol)

	validate := validation.New()

	app := http.NewServer(&logger, usersRepo, exercisesRepo, trainingsRepo, validate)

	app.AddHandlers()

	app.Run(":" + os.Getenv("PORT"))
}
