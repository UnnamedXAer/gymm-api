package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// func TestMain(m *testing.M) {
// 	return
// 	ValidateTestEnv()
// 	logger := zerolog.New(os.Stdout)
// 	db, err := GetDatabase(os.Getenv("MONGO_URI"))
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	defer DisconnectDB(db)
// 	usersCollection := db.Collection("users")
// 	trainingsCollection := db.Collection("trainings")
// 	repoUsers := users.NewRepository(&logger, usersCollection)
// 	repoTrainings := trainings.NewRepository(&logger, trainingsCollection)

// 	user, err := repoUsers.GetUserByID("600da80ec5b71e2b2a712a80")

// 	if err != nil {
// 		logger.Error().Msg(err.Error())
// 	}
// 	logger.Debug().Msgf("%v", user)

// 	training, err := repoTrainings.StartTraining("600da80ec5b71e2b2a712a80", time.Now())

// 	if err != nil {
// 		logger.Error().Msg(err.Error())
// 	}
// 	logger.Debug().Msgf("%v", training)
// }

func ValidateTestEnv() error {
	err := godotenv.Load("../.test.env")
	if err != nil {
		err := godotenv.Load("../../.test.env")
		if err != nil {
			err := godotenv.Load(".test.env")
			if err != nil {
				return (err)
			}
		}
	}
	if os.Getenv("ENV") != "TEST" {
		return fmt.Errorf("wrong env, wanted '%s', got '%s'", "TEST", os.Getenv("ENV"))
	}
	return nil
}

func DisconnectDB(db *mongo.Database) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	err := db.Client().Disconnect(ctx)
	if err != nil {
		log.Printf("db disconnect error: %v", err)
		return
	}
	log.Println("db client disconnected")
}

func GetDatabase(uri string) (*mongo.Database, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	dbName, ok := os.LookupEnv("DB_NAME")
	if ok == false {
		return nil, errors.New("db name not specified")
	}
	log.Println("DB name: ", dbName)
	db := client.Database(dbName)
	return db, nil
}
