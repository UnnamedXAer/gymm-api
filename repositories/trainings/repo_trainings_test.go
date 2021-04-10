package trainings

import (
	"context"
	"math/rand"
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
)

var (
	trainingRepo *TrainingRepository
	userRepo     *users.UserRepository
	// db                    *mongo.Database
	trainingdata          trainingData
	mockedUser            entities.User
	mockedStartedTraining entities.Training
	mockedStartedExercise entities.TrainingExercise
	mockedSet             entities.TrainingSet
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
	gotTraining, err := trainingRepo.StartTraining(trainingdata.UserID.Hex(), trainingdata.StartTime)
	mockedStartedTraining = gotTraining
	if err != nil {
		t.Errorf("expect to start training, got error: %v", err)
		return
	}

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

func TestGetStartedTrainings(t *testing.T) {
	if mockedStartedTraining.StartTime.IsZero() {
		t.Run("create new started training by 'TestStartTraining'", TestStartTraining)
	}

	gotTrainings, err := trainingRepo.GetUserTrainings(mockedStartedTraining.UserID, true)
	if err != nil {
		t.Errorf("expect to get started training, got error: %v", err)
		return
	}

	if len(gotTrainings) != 1 {
		t.Errorf("expect to get one started traning, got: %d", len(gotTrainings))
		return
	}

	gotTraining := gotTrainings[0]
	if testhelpers.TimesEqual(gotTraining.StartTime, trainingdata.StartTime) == false {
		t.Errorf("expect training start time to be: %s, got %s",
			trainingdata.StartTime,
			gotTraining.StartTime)
	}

	if gotTraining.ID == "" {
		t.Errorf("expected 'ID' to NOT be zero value, got %q", gotTraining.ID)
	}

	if gotTraining.EndTime.IsZero() == false {
		t.Errorf("expected 'EndTime' to be zero value, got %q", gotTraining.EndTime)
	}

	if len(gotTraining.Exercises) > 0 {
		t.Errorf("expected 'Exercises' to be empty, got %v", gotTraining.Exercises)
	}

	if gotTraining.Comment != "" {
		t.Errorf("expected 'Comment' to be empty, got %q", gotTraining.Comment)
	}
}

func TestAddExercise(t *testing.T) {
	if mockedStartedTraining.StartTime.IsZero() {
		t.Run("create new started training by 'TestStartTraining'", TestStartTraining)
	}

	now := time.Now()
	exId := "6070007dac9cb6e543aba500" // @todo: exId from db
	mockedStartedExercise.StartTime = now
	mockedStartedExercise.ExerciseID = exId
	var te entities.TrainingExercise
	te, err := trainingRepo.AddExercise(mockedStartedTraining.ID, mockedStartedExercise)
	if err != nil {
		t.Errorf("expect to add exercise, got error: %v", err)
		return
	}

	if !testhelpers.TimesEqual(te.StartTime, now) {
		t.Errorf("expect exercise start time to be: %s, got %s", now, te.StartTime)
	}

	if te.ID == "" {
		t.Error("expect 'ID' not to be empty", te.ID)
	}

	if te.ExerciseID != exId {
		t.Errorf("expect 'ExerciseID' to be %q, got %q", exId, te.ExerciseID)
	}

	if !te.EndTime.IsZero() {
		t.Errorf("expected 'EndTime' to be zero value, got %q", te.EndTime)
	}

	if te.Comment != "" {
		t.Errorf("expected 'Comment' to be empty, got %q", te.Comment)
	}

	if len(te.Sets) > 0 {
		t.Errorf("expected 'Sets' to be empty, got %v", te.Sets)
	}

	mockedStartedExercise = te
}

func TestAddSet(t *testing.T) {
	if mockedStartedExercise.StartTime.IsZero() {
		t.Run("create new started exercise by 'TestAddExercise'", TestAddExercise)
	}

	now := time.Now()
	reps := rand.New(rand.NewSource(time.Now().Unix())).Intn(30)
	mockedSet.Time = now
	mockedSet.Reps = reps
	var ts entities.TrainingSet
	ts, err := trainingRepo.AddSet(mockedStartedExercise.ID, mockedSet)
	if err != nil {
		t.Errorf("expect to add set, got error: %v", err)
		return
	}

	if ts.ID == "" {
		t.Errorf("expect 'ID' not to be empty, got %q", ts.ID)
	}

	if !testhelpers.TimesEqual(ts.Time, now) {
		t.Errorf("expect set time to be: %s, got %s", now, ts.Time)
	}

	if ts.Reps != reps {
		t.Errorf("expect reps to be %d, got %d", reps, ts.Reps)
	}

	mockedSet = ts
}

func TestGetTrainingExercises(t *testing.T) {
	if mockedSet.Time.IsZero() {
		t.Run("create new started exercise by 'TestAddSet'", TestAddSet)
	}

	now := time.Now()
	reps := 12
	mockedSet.Time = now
	mockedSet.Reps = reps
	var te []entities.TrainingExercise
	te, err := trainingRepo.GetTrainingExercises(mockedStartedTraining.ID)
	if err != nil {
		t.Errorf("expect to get training exercises, got error: %v", err)
		return
	}

	if (len(te)) == 0 {
		t.Errorf("expect to get at least one exercises for training %q", mockedStartedTraining.ID)
		return
	}

	var testExercise entities.TrainingExercise
	for _, ex := range te {
		if mockedStartedExercise.ID == ex.ID {
			testExercise = ex
			break
		}
	}

	if testExercise.ID == "" {
		t.Errorf("expect exercise %q to be among training %q exercises", mockedStartedExercise.ID, mockedStartedTraining.ID)
		return
	}

	var testSet entities.TrainingSet
	for _, set := range testExercise.Sets {
		if mockedSet.ID == set.ID {
			testSet = set
			break
		}
	}

	if testSet.ID == "" {
		t.Errorf("expect set %q to be among sets of exercise %q", mockedSet.ID, testExercise.ID)
	}

}

func TestEndExercise(t *testing.T) {
	if mockedStartedExercise.StartTime.IsZero() {
		t.Run("create new started exercise by 'TestAddExercise'", TestAddExercise)
	}

	now := time.Now()
	var te entities.TrainingExercise
	te, err := trainingRepo.EndExercise(mockedStartedExercise.ID, now)
	if err != nil {
		t.Errorf("expect to end exercise, got error: %v", err)
	}

	if testhelpers.TimesEqual(te.EndTime, now) {
		t.Errorf("expect exercise end time to be: %s, got %s", now, te.EndTime)
	}
}

func TestEndTraining(t *testing.T) {
	if mockedStartedTraining.StartTime.IsZero() {
		t.Run("create new started training by 'TestStartTraining'", TestStartTraining)
	}

	now := time.Now()
	tr, err := trainingRepo.EndTraining(mockedStartedTraining.ID, now)
	if err != nil {
		t.Errorf("expected to end training (%s), got error: %v", mockedStartedTraining.ID, err)
		return
	}

	if !testhelpers.TimesEqual(tr.EndTime, now) {
		t.Errorf("expected to training end time be: %s, got %s", now, tr.EndTime)
	}
}

func TestGetUserTrainings(t *testing.T) {
	if mockedStartedTraining.StartTime.IsZero() {
		t.Run("create new started training by 'TestStartTraining'", TestStartTraining)
	}

	var tr []entities.Training
	tr, err := trainingRepo.GetUserTrainings(mockedStartedTraining.UserID, false)
	if err != nil {
		t.Errorf("expected to get trainings for user %q, got error: %v", mockedStartedTraining.UserID, err)
		return
	}

	if len(tr) == 0 {
		t.Errorf("expect to get at least one training for user %q", mockedStartedTraining.UserID)
		return
	}

	var testTraining entities.Training
	for _, training := range tr {
		if mockedStartedTraining.ID == training.ID {
			testTraining = training
			break
		}
	}

	if testTraining.ID == "" {
		t.Errorf("expect training %q to be among trainings for user %q", mockedStartedTraining.ID, mockedStartedTraining.UserID)
	}
}
