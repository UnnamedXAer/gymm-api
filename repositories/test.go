package repositories

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/unnamedxaer/gymm-api/entities"
	"github.com/unnamedxaer/gymm-api/usecases"
	"go.mongodb.org/mongo-driver/mongo"
)

// EnsureTestEnv tries to load the '.test.env' file end ensure the 'ENV' is set to 'TEST'
func EnsureTestEnv() {
	err := godotenv.Load("../.test.env")
	if err != nil {
		err := godotenv.Load("../../.test.env")
		if err != nil {
			err := godotenv.Load(".test.env")
			if err != nil {
				panic(err)
			}
		}
	}
	if os.Getenv("ENV") != "test" {
		panic(fmt.Errorf("wrong env, wanted '%s', got '%s'", "test", os.Getenv("ENV")))
	}
}

// DisconnectDB disconnects from given db
func DisconnectDB(l *zerolog.Logger, db *mongo.Database) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	err := db.Client().Disconnect(ctx)
	if err != nil {
		l.Info().Msgf("db '%s' disconnect error: %v", db.Name(), err)
		return
	}
	l.Info().Msgf("db '%s' disconnected :", db.Name())
}

// InsertMockUser inserts mocked user to repository with use of the repo functionality
func InsertMockUser(ur usecases.UserRepo) (entities.User, error) {
	if os.Getenv("ENV") == "test" {
		panic(fmt.Errorf("wrong env, NOT wanted 'test', got '%s'", os.Getenv("ENV")))
	}

	return ur.CreateUser(
		"John Silver",
		"johnsilver@email.com",
		[]byte("TheSecretestPasswordEver123$%^"),
	)
}

func StartMockTraining(tr usecases.TrainingRepo) (entities.Training, error) {
	return tr.StartTraining()
}

func TimesEqual(t1, t2 time.Time) bool {

	return t1.Year() == t2.Year() &&
		t1.Month() == t2.Month() &&
		t1.Day() == t2.Day() &&
		t1.Hour() == t2.Hour() &&
		t1.Minute() == t2.Minute() &&
		t1.Second() == t2.Second()
}
