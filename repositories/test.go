package repositories

import (
	"context"
	"fmt"
	"os"
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
	if os.Getenv("ENV") != "TEST" {
		panic(fmt.Errorf("wrong env, wanted '%s', got '%s'", "TEST", os.Getenv("ENV")))
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
