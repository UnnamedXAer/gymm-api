package exercises

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/unnamedxaer/gymm-api/entities"
	"github.com/unnamedxaer/gymm-api/mocks"
	"github.com/unnamedxaer/gymm-api/repositories"
	"github.com/unnamedxaer/gymm-api/testhelpers"
	"github.com/unnamedxaer/gymm-api/usecases"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	exerciseColName = "exercises"
)

var (
	exerciseRepo   usecases.ExerciseRepo
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
	want := mockedExercise
	ex, err := exerciseRepo.GetExerciseByID(want.ID)
	if err != nil {
		t.Error(err)
		return
	}

	if ex.ID != want.ID {
		t.Errorf("want 'ID' to be %q, got %q", want.ID, ex.ID)
	}

	if ex.Name != want.Name {
		t.Errorf("want 'Name' to be %q, got %q", want.Name, ex.Name)
	}

	if ex.Description != want.Description {
		t.Errorf("want 'Description' to be %q, got %q", want.Description, ex.Description)
	}

	if !testhelpers.TimesEqual(ex.CreatedAt, want.CreatedAt) {
		t.Errorf("want 'CreateAt' to be %s, got %s", want.CreatedAt, ex.CreatedAt)
	}

	if ex.CreatedBy != want.CreatedBy {
		t.Errorf("want 'CreateBy' to be %q, got %q", want.CreatedBy, ex.CreatedBy)
	}

	if ex.SetUnit != want.SetUnit {
		t.Errorf("want 'SetUnit' to be Time (%d), got %d", want.SetUnit, ex.SetUnit)
	}
	mockedExercise = *ex
}

func TestCreateExercise(t *testing.T) {
	want := mockedExercise
	want.Name += fmt.Sprintf("-> %d", time.Now().UnixNano())
	ex, err := exerciseRepo.CreateExercise(want.Name, want.Description, want.SetUnit, want.CreatedBy)
	if err != nil {
		t.Error(err)
		return
	}

	if ex.ID == "" {
		t.Error("want 'ID' not to be zero value")
		return
	}

	if ex.Name != want.Name {
		t.Errorf("want 'Name' to be %q, got %q", want.Name, ex.Name)
	}

	if ex.Description != want.Description {
		t.Errorf("want 'Description' to be %q, got %q", want.Description, ex.Description)
	}

	if ex.CreatedAt.IsZero() {
		t.Error("want 'CreateAt' NOT to be zero value")
	}

	if ex.CreatedBy != want.CreatedBy {
		t.Errorf("want 'CreateBy' to be %q, got %q", want.CreatedBy, ex.CreatedBy)
	}

	if ex.SetUnit != want.SetUnit {
		t.Errorf("want 'SetUnit' to be Time (%d), got %d", want.SetUnit, ex.SetUnit)
	}
	mockedExercise = *ex
}

func TestUpdateExercise(t *testing.T) {
	want := mockedExercise
	want.Description += "\n-> updated at " + time.Now().String()
	want.SetUnit = entities.Time
	ex, err := exerciseRepo.UpdateExercise(&want)
	if err != nil {
		t.Error(err)
		return
	}

	if ex.ID != want.ID {
		t.Errorf("want 'ID' to be %q, got %q", want.ID, ex.ID)
	}

	if ex.Name != want.Name {
		t.Errorf("want 'Name' to be %q, got %q", want.Name, ex.Name)
	}

	if ex.Description != want.Description {
		t.Errorf("want 'Description' to be %q, got %q", want.Description, ex.Description)
	}

	if !testhelpers.TimesEqual(ex.CreatedAt, want.CreatedAt) {
		t.Errorf("want 'CreateAt' to be %s, got %s", want.CreatedAt, ex.CreatedAt)
	}

	if ex.CreatedBy != want.CreatedBy {
		t.Errorf("want 'CreateBy' to be %q, got %q", want.CreatedBy, ex.CreatedBy)
	}

	if ex.SetUnit != want.SetUnit {
		t.Errorf("want 'SetUnit' to be Time (%d), got %d", want.SetUnit, ex.SetUnit)
	}
	mockedExercise = *ex
}

func TestUpdateExerciseOneProp(t *testing.T) {
	input := entities.Exercise{
		ID:      mockedExercise.ID,
		SetUnit: entities.Time,
	}
	want := mockedExercise
	want.SetUnit = input.SetUnit
	ex, err := exerciseRepo.UpdateExercise(&input)
	if err != nil {
		t.Error(err)
		return
	}

	if ex.ID != want.ID {
		t.Errorf("want 'ID' to be %q, got %q", want.ID, ex.ID)
	}

	if ex.Name != want.Name {
		t.Errorf("want 'Name' to be %q, got %q", want.Name, ex.Name)
	}

	if ex.Description != want.Description {
		t.Errorf("want 'Description' to be %q, got %q", want.Description, ex.Description)
	}

	if !testhelpers.TimesEqual(ex.CreatedAt, want.CreatedAt) {
		t.Errorf("want 'CreateAt' to be %s, got %s", want.CreatedAt, ex.CreatedAt)
	}

	if ex.CreatedBy != want.CreatedBy {
		t.Errorf("want 'CreateBy' to be %q, got %q", want.CreatedBy, ex.CreatedBy)
	}

	if ex.SetUnit != want.SetUnit {
		t.Errorf("want 'SetUnit' to be Time (%d), got %d", want.SetUnit, ex.SetUnit)
	}
	mockedExercise = *ex
}
