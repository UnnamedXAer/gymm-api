package trainings

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/unnamedxaer/gymm-api/entities"
	"github.com/unnamedxaer/gymm-api/repositories"
	"github.com/unnamedxaer/gymm-api/repositories/users"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	trainingRepo          *TrainingRepository
	userRepo              *users.UserRepository
	db                    *mongo.Database
	trainingdata          trainingData
	mockedUser            entities.User
	mockedStartedTraining entities.Training
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
	// _, err = trainingsCol.DeleteMany(nil, bson.D{})
	// if err != nil {
	// 	panic(err)
	// }

	trainingRepo = NewRepository(&logger, trainingsCol)
	userRepo = users.NewRepository(&logger, usersCol)

	res, err := usersCol.DeleteMany(context.TODO(), bson.D{})
	if err != nil {
		logger.Err(err).Msgf("%d", res.DeletedCount)
		panic(err)
	}
	mockedUser, err = repositories.InsertMockUser(userRepo)
	if err != nil {
		panic(err)
	}
	uOID, err := primitive.ObjectIDFromHex(mockedUser.ID)
	if err != nil {
		panic(err)
	}

	trainingdata = trainingData{
		UserID:    uOID,
		StartTime: time.Now().UTC(),
	}

	code := m.Run()
	os.Exit(code)
}

func TestStartTraining(t *testing.T) {
	// clearCollection(t)
	gotTraining, err := trainingRepo.StartTraining(trainingdata.UserID.Hex(), trainingdata.StartTime)
	mockedStartedTraining = gotTraining
	if err != nil {
		t.Fatal(err)
	}
	if repositories.TimesEqual(gotTraining.StartTime, trainingdata.StartTime) == false {
		t.Fatalf("Expect to get started training base on data: %v, got: %v",
			trainingdata,
			gotTraining)
	}

	if gotTraining.ID == "" {
		t.Fatalf("Expected 'ID' to NOT be zero value, got %q", gotTraining.ID)
	}

	if gotTraining.EndTime.IsZero() == false {
		t.Fatalf("Expected 'EndTime' to be zero value, got %q", gotTraining.EndTime)
	}

	if len(gotTraining.Exercises) > 0 {
		t.Fatalf("Expected 'Exercises' to be empty, got %v exercises", gotTraining.Exercises)
	}

	if gotTraining.Comment != "" {
		t.Fatalf("Expected 'Comment' to be empty, got %q", gotTraining.Comment)
	}
}

func TestGetStartedTraining(t *testing.T) {
	if mockedStartedTraining.StartTime.IsZero() {
		t.Run("create new started training by 'TestStartTraining'", TestStartTraining)
	}

	gotTraining, err := trainingRepo.GetStartedTraining((mockedStartedTraining.UserID.Hex())
	if err != nil {
		t.Fatal(err)
	}
	if repositories.TimesEqual(gotTraining.StartTime, trainingdata.StartTime) == false {
		t.Fatalf("Expect to get started training base on data: %v, got: %v",
			trainingdata,
			gotTraining)
	}

	if gotTraining.ID == "" {
		t.Fatalf("Expected 'ID' to NOT be zero value, got %q", gotTraining.ID)
	}

	if gotTraining.EndTime.IsZero() == false {
		t.Fatalf("Expected 'EndTime' to be zero value, got %q", gotTraining.EndTime)
	}

	if len(gotTraining.Exercises) > 0 {
		t.Fatalf("Expected 'Exercises' to be empty, got %v exercises", gotTraining.Exercises)
	}

	if gotTraining.Comment != "" {
		t.Fatalf("Expected 'Comment' to be empty, got %q", gotTraining.Comment)
	}
}

func clearCollection(t *testing.T) {
	_, err := trainingRepo.col.DeleteMany(context.TODO(), bson.D{})
	if err != nil {
		t.Fatal(err)
	}
}

func insertMockTraining(t *testing.T) {
	_, err := trainingRepo.col.InsertOne(context.TODO(), trainingdata)
	if err != nil {
		t.Fatal(err)
	}
}
