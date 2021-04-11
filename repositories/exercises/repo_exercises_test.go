package exercises

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/unnamedxaer/gymm-api/entities"
	"github.com/unnamedxaer/gymm-api/mocks"
	"github.com/unnamedxaer/gymm-api/repositories"
	"github.com/unnamedxaer/gymm-api/testhelpers"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	exerciseColName = "exercises"
)

var (
	exerciseRepo   *ExerciseRepository
	mockedExercise entities.Exercise
)

func TestMain(m *testing.M) {

	testhelpers.EnsureTestEnv()
	logger := zerolog.New(os.Stdout)

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatalln("environment variable 'DB_NAME' is not set")
	}
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatalln("environment variable 'MONGO_URI' is not set")
	}
	db, err := repositories.GetDatabase(&logger, mongoURI, dbName)
	if err != nil {
		log.Fatalln(err)
	}

	err = repositories.CreateCollections(&logger, db)
	if err != nil {
		log.Fatalln(err)
	}
	defer testhelpers.DisconnectDB(&logger, db)

	exercisesCol := db.Collection(exerciseColName)
	res, err := exercisesCol.DeleteMany(context.TODO(), bson.D{})
	if err != nil {
		log.Fatalln(err)
	}

	exerciseRepo = NewRepository(&logger, exercisesCol)

	if err != nil {
		logger.Err(err).Msgf("%d", res.DeletedCount)
		panic(err)
	}
	me, err := mocks.InsertMockExercise(exerciseRepo)
	if err != nil {
		panic(err)
	}
	mockedExercise = *me

	os.Exit(m.Run())
}

func TestGetExerciseByID(t *testing.T) {
	ex, err := exerciseRepo.GetExerciseByID(mockedExercise.ID)
	if err != nil {
		t.Error(err)
		return
	}

	if ex.ID == "" {
		t.Error("expect 'ID' not to be zero value")
		return
	}

	if ex.Name != mockedExercise.Name {
		t.Errorf("expect 'Name' to be %q, got %q", mockedExercise.Name, ex.Name)
	}

	if ex.Description != mockedExercise.Description {
		t.Errorf("expect 'Description' to be %q, got %q", mockedExercise.Description, ex.Description)
	}

	if ex.CreatedAt.IsZero() {
		t.Error("expect 'CreateAt' NOT to be zero value")
	}

	if ex.CreatedBy == "" {
		t.Error("expect 'CreateBy' NOT to be zero value")
	}

	if ex.SetUnit == 0 {
		t.Error("expect 'SetUnit' NOT to be zero value")
	}
}
