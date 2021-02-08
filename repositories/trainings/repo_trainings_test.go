package trainings

import (
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/unnamedxaer/gymm-api/repositories"
	"github.com/unnamedxaer/gymm-api/repositories/users"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	tr *TrainingRepository
	ur *users.UserRepository
	db *mongo.Database
	t  trainingData
)

const (
	trCollName = "trainings"
	uCollName  = "users"
)

func TestMain(m *testing.M) {
	repositories.EnsureTestEnv()
	logger := zerolog.New(os.Stdout)

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

	err = repositories.CreateCollections(&logger, db)
	if err != nil {
		panic(err)
	}
	defer repositories.DisconnectDB(&logger, db)

	trainingsCol := db.Collection(trCollName)
	usersCol := db.Collection(uCollName)
	_, err = trainingsCol.DeleteMany(nil, bson.D{})
	if err != nil {
		panic(err)
	}

	tr = NewRepository(&logger, trainingsCol)
	ur = users.NewRepository(&logger, usersCol)

	u, err := repositories.InsertMockUser(ur)
	if err != nil {
		panic(err)
	}
	uOID, err := primitive.ObjectIDFromHex(u.ID)
	if err != nil {
		panic(err)
	}

	t = trainingData{
		UserID:    uOID,
		StartTime: time.Now().UTC(),
	}

	code := m.Run()
	os.Exit(code)
}
