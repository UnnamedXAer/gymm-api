package testhelpers

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
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
	if strings.ToLower(os.Getenv("ENV")) != "test" {
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

func TimesEqual(t1, t2 time.Time) bool {

	return t1.Year() == t2.Year() &&
		t1.Month() == t2.Month() &&
		t1.Day() == t2.Day() &&
		t1.Hour() == t2.Hour() &&
		t1.Minute() == t2.Minute() &&
		t1.Second() == t2.Second()
}
