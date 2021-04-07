package trainings

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/unnamedxaer/gymm-api/entities"
	"github.com/unnamedxaer/gymm-api/mocks"
	"github.com/unnamedxaer/gymm-api/repositories"
	"github.com/unnamedxaer/gymm-api/repositories/users"
	"github.com/unnamedxaer/gymm-api/testhelpers"
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
	testhelpers.EnsureTestEnv()
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
	defer testhelpers.DisconnectDB(&logger, db)

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
	mockedUser, err = mocks.InsertMockUser(userRepo)
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
	if testhelpers.TimesEqual(gotTraining.StartTime, trainingdata.StartTime) == false {
		t.Errorf("Expect to get started training base on data: %v, got: %v",
			trainingdata,
			gotTraining)
		return
	}

	if gotTraining.ID == "" {
		t.Errorf("expected 'ID' to NOT be zero value, got %q", gotTraining.ID)
	}

	if gotTraining.EndTime.IsZero() == false {
		t.Errorf("expected 'EndTime' to be zero value, got %q", gotTraining.EndTime)
	}

	if len(gotTraining.Exercises) > 0 {
		t.Errorf("expected 'Exercises' to be empty, got %v exercises", gotTraining.Exercises)
	}

	if gotTraining.Comment != "" {
		t.Errorf("expected 'Comment' to be empty, got %q", gotTraining.Comment)
	}
}

func TestGetStartedTrainings(t *testing.T) {
	if mockedStartedTraining.StartTime.IsZero() {
		t.Run("create new started training by 'TestStartTraining'", TestStartTraining)
	}

	gotTrainings, err := trainingRepo.GetStartedTrainings(mockedStartedTraining.UserID)
	if err != nil {
		t.Errorf("expect to get one training, got: %v", err)
		return
	}

	if len(gotTrainings) != 1 {
		t.Errorf("expect to get one started traning, got: %d", len(gotTrainings))
		return
	}

	gotTraining := gotTrainings[0]
	if testhelpers.TimesEqual(gotTraining.StartTime, trainingdata.StartTime) == false {
		t.Errorf("expect to get started training base on data: %v, got: %v",
			trainingdata,
			gotTraining)
		return
	}

	if gotTraining.ID == "" {
		t.Errorf("expected 'ID' to NOT be zero value, got %q", gotTraining.ID)
	}

	if gotTraining.EndTime.IsZero() == false {
		t.Errorf("expected 'EndTime' to be zero value, got %q", gotTraining.EndTime)
	}

	if len(gotTraining.Exercises) > 0 {
		t.Errorf("expected 'Exercises' to be empty, got %v exercises", gotTraining.Exercises)
	}

	if gotTraining.Comment != "" {
		t.Errorf("expected 'Comment' to be empty, got %q", gotTraining.Comment)
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
